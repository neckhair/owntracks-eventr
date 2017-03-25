# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  # config.vm.synced_folder "../data", "/vagrant_data"
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "512"
  end

  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y mosquitto
  SHELL
end
