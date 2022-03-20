import sys

unit_size = 10
count = 10000

if len(sys.argv) == 3:
    unit_size, count = int(sys.argv[1]), int(sys.argv[2])

if unit_size == 10:
    value = "1010101010"
elif unit_size == 20:
    value = "20202020202020202020"
elif unit_size == 50:
    value = "50505050505050505050505050505050505050505050505050"

redis_commands = ""
for i in range(count):
    redis_commands += f"SET key{i} {value}\n"
    
print(redis_commands)
