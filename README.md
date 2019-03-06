# 日志流量统计

安装docker

    sudo wget -qO- https://get.docker.com/ | sh
		
安装influxdb

    docker pull influxdb
    
    docker run --name my_influxdb influxdb
    
    docker exec -it influxdb bash
    
    进入： influxdb
    
    创建数据库：create database log
    
    进入数据库：use log

		
安装grafana	

    docker run -d -p 3000:3000 grafana/grafana --name my_grafana
    
    访问：http://Ip:3000
    
    默认用户/密码： admin/admin

