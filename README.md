# RClip
RClip is a headless clipboard.

## Installation
```
git clone https://github.com/NightWolf007/rclip.git
cd rclip
dep ensure
sudo go build -o /usr/bin/rclip
sudo mkdir /etc/rclip
sudo cp rclipd.sample.yaml /etc/rclip/rclipd.yaml
sudo cp rclipd.service /lib/systemd/system/rclipd.service
sudo systemctl daemon-reload
```

Now you can manage it through systemd:
```
sudo systemctl enable rclipd
sudo systemcll start rclipd
```
