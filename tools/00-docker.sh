#!/bin/bashsudo apt-get install -y \

apt-transport-https \
ca-certificates \
curl \
software-properties-common curl -fsSL https://yum.dockerproject.org/gpg | sudo apt-key add - sudo add-apt-repository \
"deb https://apt.dockerproject.org/repo/ \
ubuntu-$(lsb_release -cs) \
main" sudo apt-get update

# Добавляю текущего пользователя в группу докеров, чтобы не было необходимости использовать sudo при запуске докера
sudo apt-get -y install docker-engine 
sudo usermod -aG docker $(whoami)