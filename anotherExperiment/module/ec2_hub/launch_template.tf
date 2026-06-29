resource "aws_launch_template" "example" {
  name_prefix   = "scroller-template-"
  image_id      = var.ami_id
  instance_type = "t3.micro"

  key_name = var.keys

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.security.id]
    subnet_id = data.aws_subnet.public_subnet.id
  }
  iam_instance_profile {
    name = var.scroller_profile
  }
  block_device_mappings {
    #device_name = "/dev/xvda"
    device_name = "/dev/sda1"
    ebs {
      volume_size           = 16
      volume_type           = "gp3"
      delete_on_termination = true
    }
  }

  instance_market_options {
    market_type = "spot"

    spot_options {
      spot_instance_type             = "one-time"
      instance_interruption_behavior = "terminate"
    }
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = "my-instance"
    }
  }

  user_data = base64encode(<<-EOF
            #!/bin/bash
            echo "Hello from launch template"
            sudo bash -c '
            Xvfb :1 -screen 0 1366x768x24 &
            export DISPLAY=:1
            echo "DISPLAY=:1" >> /etc/environment
            until xdpyinfo -display :1 >/dev/null 2>&1; do sleep 1; done
            echo "Display ready!"
            startxfce4 &
            until pgrep -x "xfce4-panel" >/dev/null 2>&1; do sleep 1; done
            echo "XFCE ready!"
            '
            EOF
  )

  monitoring {
    enabled = true
  }
}