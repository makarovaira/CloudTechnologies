terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
    telegram = {
      source  = "yi-jiayu/telegram"
      version = "0.3.1"
    }
  }
  required_version = ">= 0.13"
}

data "archive_file" "content" {
  type        = "zip"
  source_dir  = "src"
  output_path = "build/content.zip"
}

data "http" "set_webhook_tg" {
  url = "https://api.telegram.org/bot${var.tg_bot_key}/setWebhook?url=https://functions.yandexcloud.net/${yandex_function.cheatsheet_itis_function.id}"
  retry {
    attempts = 5
    min_delay_ms = 1000
    max_delay_ms = 1000
  }
}

variable "cloud_id" {
  type        = string
  description = "ID облака"
}

variable "folder_id" {
  type        = string
  description = "ID каталога"
}

variable "tg_bot_key" {
  type        = string
  description = "Токен Telegram Bot API"
}

variable "sa_key_file_path" {
  type        = string
  description = "Путь ключа для провайдера Yandex.Cloud"
  default     = "~/.yc-keys/key.json"
}

variable "bucket_name" {
  type        = string
  description = "Название бакета для инструкции YandexGPT"
  default = "gpt-prompt"
}

variable "bucket_object_key" {
  type        = string
  description = "Инструкция для YandexGPT"
  default     = "prompt.txt"
}

provider "yandex" {
  cloud_id                 = var.cloud_id
  folder_id                = var.folder_id
  zone                     = "ru-central1-a"
  service_account_key_file = pathexpand(var.sa_key_file_path)
}

provider "telegram" {
  bot_token = var.tg_bot_key
}

resource "yandex_function" "cheatsheet_itis_function" {
  name               = "cheatsheet-itis-function"
  entrypoint         = "main.Handler"
  memory             = "128"
  runtime            = "golang121"
  service_account_id = yandex_iam_service_account.sa.id
  user_hash          = data.archive_file.content.output_sha512
  execution_timeout  = "60"
  environment = {
    TELEGRAM_BOT_TOKEN = var.tg_bot_key
    FOLDER_ID          = var.folder_id
    BUCKET_OBJECT_KEY  = var.bucket_object_key
  }
  content {
    zip_filename = data.archive_file.content.output_path
  }
  mounts {
    name = "mnt"
    mode = "ro"
    object_storage {
      bucket = yandex_storage_bucket.exam_solver_tg_bot_bucket.bucket
    }
  }
}

resource "yandex_iam_service_account" "sa" {
  name = var.bucket_name
}

resource "yandex_function_iam_binding" "exam_solver_tg_bot_iam" {
  function_id = yandex_function.cheatsheet_itis_function.id
  role        = "functions.functionInvoker"
  members = [
    "system:allUsers",
  ]
}

resource "yandex_resourcemanager_folder_iam_member" "sa_exam_solver_tg_bot_ai_vision_iam" {
  folder_id = var.folder_id
  role      = "ai.vision.user"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "sa_exam_solver_tg_bot_ai_language_models_iam" {
  folder_id = var.folder_id
  role      = "ai.languageModels.user"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "sa_exam_solver_tg_bot_storage_viewer_iam" {
  folder_id = var.folder_id
  role      = "storage.admin"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

resource "telegram_bot_webhook" "exam_solver_tg_bot_webhook" {
  url = "https://functions.yandexcloud.net/${yandex_function.cheatsheet_itis_function.id}"
}

resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.sa.id
  description        = "static access key for object storage"
}

resource "yandex_storage_bucket" "exam_solver_tg_bot_bucket" {
  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key

  bucket = var.bucket_name
}

resource "yandex_storage_object" "yandexgpt_instruction" {
  bucket = yandex_storage_bucket.exam_solver_tg_bot_bucket.id
  key    = var.bucket_object_key
  source = "prompt.txt"
}
