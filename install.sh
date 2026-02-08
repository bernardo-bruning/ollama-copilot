#/bin/bash

go build
mkdir -p $HOME/.local/bin
cp ollama-copilot $HOME/.local/bin/ollama-copilot

echo "Install in systemd? [y/N]"
read INSTALL_SYSTEMD
if [[ $INSTALL_SYSTEMD == "y" ]]
then
  echo "instalando systemctl"
  cp ollama-copilot.service $HOME/.config/systemd/user/ollama-copilot.service
  systemctl --user daemon-reload
  systemctl --user enable ollama-copilot.service
  systemctl --user start ollama-copilot.service
fi

echo $INSTALL_SYSTEMD
