---
id: system-monitor
spinal: true
log: true
---

## Measure ∡ computer metrics

```py
import psutil
import crython
from time import sleep

# Every 10 seconds
@crython.job(second='*/10')
def heart_beat():
  print('Heartbeat ❤')

# At every minute of every hour at second 1
@crython.job(second=0, minute='*/1')
def calc_metrics():
  cpu = psutil.cpu_percent(interval=1)
  virt_mem = psutil.virtual_memory()
  swap_mem = psutil.virtual_memory()
  disk = psutil.disk_usage('/')

  print('CPU:', cpu, '%')
  print('Virt MEM:', virt_mem.percent, '%')
  print('Swap MEM:', swap_mem.percent, '%')
  print('Disk usage:', disk.percent, '%')

if __name__ == '__main__':
  print('System metrics Ⓜ️')
  crython.start()
  while True:
    sleep(1)
```

## Bye-bye metrics Ⓜ️
