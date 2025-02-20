# Приватное
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
  description = "Токен для доступа к Telegram Bot API"
}

# Очередь
variable "face_cut_queue" {
  type        = string
  description = "Название очереди для заданий по обрезке лиц"
  default     = "vvot20-task"
}

# Тригеры
variable "face_detection_trigger" {
  type        = string
  description = "Название триггера бакета для обнаружения лиц"
  default     = "vvot20-photo"
}

variable "face_cut_trigger" {
  type        = string
  description = "Название триггера очереди для заданий обрезки лиц"
  default     = "vvot20-task"
}

# Функции
variable "face_detection_function" {
  type        = string
  description = "Название функции для обнаружения лиц"
  default     = "vvot20-face-detection"
}

variable "face_cut_function" {
  type        = string
  description = "Название функции для обрезания лиц"
  default     = "vvot20-face-cut"
}


variable "bot_function" {
  type        = string
  description = "Название функции для Telegram-бота"
  default     = "vvot20-boot"
}

# Бакеты
variable "photos_bucket" {
  type        = string
  description = "Бакет для оригинальных фотографий"
  default     = "vvot20-photo"
}

variable "faces_bucket" {
  type        = string
  description = "Бакет для фотографий лиц"
  default     = "vvot20-faces"
}

# Шлюз
variable "api_gateway" {
  type        = string
  description = "Название API-шлюза для фотографий лиц"
  default     = "vvot20-apigw"
}

# Сервисные аккаунты
variable "sa_key_file_path" {
  type        = string
  description = "Путь для Провайдер «Yandex.Cloud Provider» чтобы искать авторизованный ключ "
  default     = "~/.yc-keys/key.json"
}

variable "sa_face_detection" {
  type        = string
  description = "Сервисный аккаунт для функции обнаружения лиц"
  default     = "sa-recognizer"
}

variable "sa_face_cut" {
  type        = string
  description = "Сервисный аккаунт для функции вырезания лиц"
  default     = "sa-cropper"
}

variable "sa_bot" {
  type        = string
  description = "Сервисный аккаунт для работы с тг ботом"
  default     = "sa-bot"
}