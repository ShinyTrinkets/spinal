---
id: system-monitor
spinal: true
log: true
---

## Measure ∡ computer metrics

```py
import psutil

print('System metrics Ⓜ️')

# Python 3.x compatible
def calc_metrics():
  cpu = psutil.cpu_percent(interval=1)
  virt_mem = psutil.virtual_memory()
  swap_mem = psutil.virtual_memory()
  disk = psutil.disk_usage('/')

  print('CPU:', cpu, '%')
  print('Virt MEM:', virt_mem.percent, '%')
  print('Swap MEM:', swap_mem.percent, '%')
  print('Disk usage:', disk.percent, '%')

calc_metrics()
```

## Bye-bye metrics Ⓜ️
