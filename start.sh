#!/bin/bash

app_name=financelog
./$app_name  ./etc/$app_name.cfg  >>./log/financelog_$(date +%Y%m%d).log 2>&1 &
