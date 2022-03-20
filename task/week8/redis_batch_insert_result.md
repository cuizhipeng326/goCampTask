## 环境：
redis_version:6.0.8



## 插入10000条value大小为10字节的数据前，redis memory情况：

```json
used_memory:4009896
used_memory_human:3.82M
used_memory_rss:14127104
used_memory_rss_human:13.47M
used_memory_peak:8076784
used_memory_peak_human:7.70M
used_memory_peak_perc:49.65%
used_memory_overhead:3218722
used_memory_startup:526760
used_memory_dataset:791174
used_memory_dataset_perc:22.71%
allocator_allocated:4123960
allocator_active:5156864
allocator_resident:9359360
total_system_memory:16792739840
total_system_memory_human:15.64G
used_memory_lua:47104
used_memory_lua_human:46.00K
used_memory_scripts:1512
used_memory_scripts_human:1.48K
number_of_cached_scripts:5
maxmemory:12884901888
maxmemory_human:12.00G
maxmemory_policy:volatile-lru
allocator_frag_ratio:1.25
allocator_frag_bytes:1032904
allocator_rss_ratio:1.81
allocator_rss_bytes:4202496
rss_overhead_ratio:1.51
rss_overhead_bytes:4767744
mem_fragmentation_ratio:3.52
mem_fragmentation_bytes:10117280
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:2683818
mem_aof_buffer:0
mem_allocator:jemalloc-5.1.0
active_defrag_running:0
lazyfree_pending_objects:0
```

## 插入10000条value大小为10字节的数据后，redis memory情况：
```json
used_memory:4652152
used_memory_human:4.44M
used_memory_rss:14598144
used_memory_rss_human:13.92M
used_memory_peak:8076784
used_memory_peak_human:7.70M
used_memory_peak_perc:57.60%
used_memory_overhead:3749762
used_memory_startup:526760
used_memory_dataset:902390
used_memory_dataset_perc:21.87%
allocator_allocated:4881104
allocator_active:5918720
allocator_resident:10121216
total_system_memory:16792739840
total_system_memory_human:15.64G
used_memory_lua:43008
used_memory_lua_human:42.00K
used_memory_scripts:1512
used_memory_scripts_human:1.48K
number_of_cached_scripts:5
maxmemory:12884901888
maxmemory_human:12.00G
maxmemory_policy:volatile-lru
allocator_frag_ratio:1.21
allocator_frag_bytes:1037616
allocator_rss_ratio:1.71
allocator_rss_bytes:4202496
rss_overhead_ratio:1.44
rss_overhead_bytes:4476928
mem_fragmentation_ratio:3.06
mem_fragmentation_bytes:9823200
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:2683818
mem_aof_buffer:0
mem_allocator:jemalloc-5.1.0
active_defrag_running:0
lazyfree_pending_objects:0
```

## 将之前插入的10000条数据的value改为20字节大小的数据后，redis memory情况：
```json
used_memory:5097832
used_memory_human:4.86M
used_memory_rss:15409152
used_memory_rss_human:14.70M
used_memory_peak:8076784
used_memory_peak_human:7.70M
used_memory_peak_perc:63.12%
used_memory_overhead:3749762
used_memory_startup:526760
used_memory_dataset:1348070
used_memory_dataset_perc:29.49%
allocator_allocated:5464472
allocator_active:6680576
allocator_resident:10145792
total_system_memory:16792739840
total_system_memory_human:15.64G
used_memory_lua:46080
used_memory_lua_human:45.00K
used_memory_scripts:1512
used_memory_scripts_human:1.48K
number_of_cached_scripts:5
maxmemory:12884901888
maxmemory_human:12.00G
maxmemory_policy:volatile-lru
allocator_frag_ratio:1.22
allocator_frag_bytes:1216104
allocator_rss_ratio:1.52
allocator_rss_bytes:3465216
rss_overhead_ratio:1.52
rss_overhead_bytes:5263360
mem_fragmentation_ratio:2.88
mem_fragmentation_bytes:10065680
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:2683818
mem_aof_buffer:0
mem_allocator:jemalloc-5.1.0
active_defrag_running:0
lazyfree_pending_objects:0
```

## 将之前插入的10000条数据的value改为50字节大小的数据后，redis memory情况：
```json
used_memory:5335920
used_memory_human:5.09M
used_memory_rss:15495168
used_memory_rss_human:14.78M
used_memory_peak:8076784
used_memory_peak_human:7.70M
used_memory_peak_perc:66.06%
used_memory_overhead:3749760
used_memory_startup:526760
used_memory_dataset:1586160
used_memory_dataset_perc:32.98%
allocator_allocated:5447264
allocator_active:6483968
allocator_resident:10280960
total_system_memory:16792739840
total_system_memory_human:15.64G
used_memory_lua:45056
used_memory_lua_human:44.00K
used_memory_scripts:1512
used_memory_scripts_human:1.48K
number_of_cached_scripts:5
maxmemory:12884901888
maxmemory_human:12.00G
maxmemory_policy:volatile-lru
allocator_frag_ratio:1.19
allocator_frag_bytes:1036704
allocator_rss_ratio:1.59
allocator_rss_bytes:3796992
rss_overhead_ratio:1.51
rss_overhead_bytes:5214208
mem_fragmentation_ratio:2.90
mem_fragmentation_bytes:10159312
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:2683816
mem_aof_buffer:0
mem_allocator:jemalloc-5.1.0
active_defrag_running:0
lazyfree_pending_objects:0
```



## 结论：

10000个10字节大小的数据占用的内存大小：4652152 - 4009896 = 642256，平均每个 key 的占用内存空间为64字节；

10000个20字节大小的数据占用的内存大小：5097832 - 4009896 = 1087936，平均每个 key 的占用内存空间为109字节；

10000个50字节大小的数据占用的内存大小：5335920 - 4009896 = 1326024，平均每个 key 的占用内存空间为133字节；

