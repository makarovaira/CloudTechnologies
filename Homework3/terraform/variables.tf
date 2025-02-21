variable "cloud_id" {
  type        = string
  description = "id облака"
}

variable "folder_id" {
  type        = string
  description = "id каталога"
}

variable "vm_user_login" {
  type = string
}

variable "zone" {
  type        = string
  description = "Yandex.Cloud Zone"
  default     = "ru-central1-d"
}