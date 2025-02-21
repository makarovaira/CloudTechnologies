resource "yandex_vpc_network" "network" {
  name = "vvot20-nextcloud-network"
}

resource "yandex_vpc_subnet" "subnet" {
  name       = "vvot20-nextcloud-subnet"
  zone       = var.zone
  v4_cidr_blocks = ["192.168.10.0/24"]
  network_id = yandex_vpc_network.network.id
}

resource "yandex_compute_disk" "boot-disk" {
  name     = "vvot20-nextcloud-boot-disk"
  type     = "network-ssd"
  image_id = data.yandex_compute_image.ubuntu.id
  size     = 20
}

resource "yandex_compute_instance" "server" {
  name        = "vvot20-nextcloud-server"
  platform_id = "standard-v3"
  hostname    = "nextcloud"

  resources {
    core_fraction = 20
    cores         = 2
    memory        = 4
  }

  boot_disk {
    disk_id = yandex_compute_disk.boot-disk.id
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.subnet.id
    nat       = true
  }

  metadata = {
    ssh-keys = "${var.vm_user_login}:${file("~/.ssh/id_rsa.pub")}"
  }
}

data "yandex_compute_image" "ubuntu" {
  family = "ubuntu-2404-lts-oslogin"
}

output "nextcloud-ip" {
  value = yandex_compute_instance.server.network_interface[0].nat_ip_address
}