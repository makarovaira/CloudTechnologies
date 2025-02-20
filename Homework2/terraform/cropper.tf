resource "yandex_storage_bucket" "faces_bucket" {
  bucket = var.faces_bucket
}

resource "yandex_function" "cropper" {
  name               = var.face_cut_function
  entrypoint         = "handler.handle"
  memory             = "128"
  runtime            = "python312"
  user_hash          = data.archive_file.cropper_source.output_sha512
  service_account_id = yandex_iam_service_account.sa_cropper.id
  environment = {
    PHOTOS_MOUNT_POINT = yandex_storage_bucket.photos_bucket.bucket
    FACES_MOUNT_POINT  = yandex_storage_bucket.faces_bucket.bucket
    ACCESS_KEY_ID      = yandex_iam_service_account_static_access_key.sa_cropper_static_key.access_key
    SECRET_ACCESS_KEY  = yandex_iam_service_account_static_access_key.sa_cropper_static_key.secret_key
  }
  content {
    zip_filename = data.archive_file.cropper_source.output_path
  }
  mounts {
    name = yandex_storage_bucket.photos_bucket.bucket
    mode = "ro"
    object_storage {
      bucket = yandex_storage_bucket.photos_bucket.bucket
    }
  }
  mounts {
    name = yandex_storage_bucket.faces_bucket.bucket
    mode = "rw"
    object_storage {
      bucket = yandex_storage_bucket.faces_bucket.bucket
    }
  }
}

resource "yandex_function_trigger" "crop_tasks_queue_trigger" {
  name = var.face_cut_trigger
  function {
    id                 = yandex_function.cropper.id
    service_account_id = yandex_iam_service_account.sa_cropper.id
  }
  message_queue {
    batch_cutoff       = "0"
    batch_size         = "1"
    queue_id           = yandex_message_queue.crop_tasks_queue.arn
    service_account_id = yandex_iam_service_account.sa_cropper.id
  }
}

data "archive_file" "cropper_source" {
  type        = "zip"
  source_dir  = "../src/cropper"
  output_path = "../build/cropper.zip"
}

resource "yandex_iam_service_account" "sa_cropper" {
  name = var.sa_face_cut
}

resource "yandex_resourcemanager_folder_iam_member" "sa_cropper_function_invoker_role" {
  folder_id = var.folder_id
  role      = "functions.functionInvoker"
  member    = "serviceAccount:${yandex_iam_service_account.sa_cropper.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "sa_cropper_storage_editor_iam" {
  folder_id = var.folder_id
  role      = "storage.editor"
  member    = "serviceAccount:${yandex_iam_service_account.sa_cropper.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "sa_cropper_ymq_reader_iam" {
  folder_id = var.folder_id
  role      = "ymq.reader"
  member    = "serviceAccount:${yandex_iam_service_account.sa_cropper.id}"
}

resource "yandex_iam_service_account_static_access_key" "sa_cropper_static_key" {
  service_account_id = yandex_iam_service_account.sa_cropper.id
}
