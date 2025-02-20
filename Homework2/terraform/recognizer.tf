resource "yandex_storage_bucket" "photos_bucket" {
  bucket = var.photos_bucket
}

resource "yandex_function" "recognizer" {
  name               = var.face_detection_function
  entrypoint         = "handler.handle"
  memory             = "256"
  runtime            = "python312"
  user_hash          = data.archive_file.recognizer_source.output_sha512
  service_account_id = yandex_iam_service_account.sa_recognizer.id
  environment = {
    MOUNT_POINT       = yandex_storage_bucket.photos_bucket.bucket
    QUEUE_URL         = yandex_message_queue.crop_tasks_queue.id
    ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.sa_recognizer_static_key.access_key
    SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.sa_recognizer_static_key.secret_key
  }
  content {
    zip_filename = data.archive_file.recognizer_source.output_path
  }
  mounts {
    name = yandex_storage_bucket.photos_bucket.bucket
    mode = "ro"
    object_storage {
      bucket = yandex_storage_bucket.photos_bucket.bucket
    }
  }
}

resource "yandex_function_trigger" "photos_bucket_trigger" {
  name = var.face_detection_trigger
  function {
    id                 = yandex_function.recognizer.id
    service_account_id = yandex_iam_service_account.sa_recognizer.id
  }
  object_storage {
    bucket_id    = yandex_storage_bucket.photos_bucket.id
    suffix       = ".jpg"
    create       = true
    batch_cutoff = "0"
  }
}

resource "yandex_message_queue" "crop_tasks_queue" {
  name       = var.face_cut_queue
  access_key = yandex_iam_service_account_static_access_key.sa_recognizer_static_key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_recognizer_static_key.secret_key
}

data "archive_file" "recognizer_source" {
  type        = "zip"
  source_dir  = "../src/recognizer"
  output_path = "../build/recognizer.zip"
}

resource "yandex_iam_service_account" "sa_recognizer" {
  name = var.sa_face_detection
}

resource "yandex_resourcemanager_folder_iam_member" "sa_recognizer_function_invoker_role" {
  folder_id = var.folder_id
  role      = "functions.functionInvoker"
  member    = "serviceAccount:${yandex_iam_service_account.sa_recognizer.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "sa_recognizer_storage_viewer_iam" {
  folder_id = var.folder_id
  role      = "storage.viewer"
  member    = "serviceAccount:${yandex_iam_service_account.sa_recognizer.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "sa_recognizer_ymq_writer_iam" {
  folder_id = var.folder_id
  role      = "ymq.writer"
  member    = "serviceAccount:${yandex_iam_service_account.sa_recognizer.id}"
}

resource "yandex_iam_service_account_static_access_key" "sa_recognizer_static_key" {
  service_account_id = yandex_iam_service_account.sa_recognizer.id
}
