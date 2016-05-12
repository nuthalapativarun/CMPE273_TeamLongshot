import bluetooth
import requests
import time

print("Scanning devices....")
list=[]
start_time=time.time()
while True:
	nearby_devices = bluetooth.discover_devices(lookup_names = True)
	if (time.time()-start_time)>30:
    	list=[]
        start_time=time.time()
	#print("found %d devices" % len(nearby_devices))
	for addr, name in nearby_devices:
		if addr not in list:   
			list.append(addr)
	 		print("  %s - %s" % (addr, name))
	 		data1={}
			request=requests.put('http://54.191.181.24:3000/profile/'+str(addr),data=data1)