

database:
  server: "localhost"
  service_name: "orcl"
  username: "system"
  password: "***"


check_items:
  - name: 数据库信息
    sql: select b.name,b.DB_UNIQUE_NAME,b.dbid,to_char(b.created,'yyyy-mm-dd hh24:mi:ss') db_created,b.database_role,b.open_mode,c.*
         from v$database b,(select v1||'_'||v2||'.'||v3 nls_lang from
          (select value v1 from nls_database_parameters where parameter='NLS_LANGUAGE') m,
          (select value v2 from nls_database_parameters where parameter='NLS_TERRITORY') n,
          (select value v3 from nls_database_parameters where parameter='NLS_CHARACTERSET') p) c

  - name: 数据库常用参数
    sql: select inst_Id,name,value from gv$parameter
         where name in ('processes','memory_max_target','memory_target','pga_aggregate_limit',
         'pga_aggregate_target','sga_max_size','sga_target','db_cache_size','shared_pool_size','java_pool_size','large_pool_size') order by 1,2

  - name: 检查表空间使用率
    sql:  SELECT
            d.status
                , d.tablespace_name
                , d.contents
                , d.extent_management
                , d.segment_space_management
                , NVL(b.allocatesize - NVL(f.freesize, 0), 0)   used_MB
                , b.allocatesize current_size_MB
                , to_char(NVL((b.allocatesize - NVL(f.freesize, 0)) / b.allocatesize * 100, 0),'990.99')||'%' pct_used
                , a.maxsize canextend_size_MB
                , to_char(NVL((b.allocatesize - NVL(f.freesize, 0)) / a.maxsize * 100, 0),'990.99')||'%' tot_pct_used
        FROM dba_tablespaces d
                , (     SELECT tablespace_name,sum(maxsize) maxsize
                        FROM (  SELECT tablespace_name, decode(autoextensible,'YES',round(sum(maxbytes)/1024/1024),round(sum(bytes)/1024/1024)) maxsize
                                        FROM dba_data_files
                                        GROUP BY tablespace_name,autoextensible
                                ) GROUP BY tablespace_name
                  ) a
                , ( SELECT tablespace_name, sum(bytes)/1024/1024 allocatesize
              from dba_data_files
              group by tablespace_name
                  ) b
                , (     SELECT tablespace_name, sum(bytes)/1024/1024 freesize
                        FROM dba_free_space
                        GROUP BY tablespace_name
                  ) f
        WHERE d.tablespace_name = a.tablespace_name(+)
        AND d.tablespace_name = b.tablespace_name(+)
        AND d.tablespace_name = f.tablespace_name(+)
        AND d.contents='PERMANENT'
        UNION ALL
        SELECT
            d.status
                , d.tablespace_name
                , d.contents
                , d.extent_management
                , d.segment_space_management
                , NVL(b.allocatesize - NVL(f.usedsize, 0), 0)   used_MB
                , b.allocatesize current_size_MB
                , to_char(NVL(NVL(f.usedsize, 0) / b.allocatesize * 100, 0),'990.99')||'%' pct_used
                , a.maxsize canextend_size_MB
                , to_char(NVL(f.usedsize,0) / a.maxsize * 100,'990.99')||'%' tot_pct_used
        FROM
            sys.dba_tablespaces d
                , (     SELECT tablespace_name,sum(maxsize) maxsize
                        FROM (  SELECT tablespace_name, decode(autoextensible,'YES',round(sum(maxbytes)/1024/1024),round(sum(bytes)/1024/1024)) maxsize
                                        FROM dba_temp_files
                                        GROUP BY tablespace_name,autoextensible
                                ) GROUP BY tablespace_name
                  ) a
          , ( select tablespace_name, sum(bytes)/1024/1024  allocatesize
              from dba_temp_files
              group by tablespace_name
            ) b
          , ( select tablespace_name, sum(bytes_cached)/1024/1024 usedsize
              from v$temp_extent_pool
              group by tablespace_name
            ) f
        WHERE d.tablespace_name = a.tablespace_name(+)
          AND d.tablespace_name = b.tablespace_name(+)
          AND d.tablespace_name = f.tablespace_name(+)
          AND d.extent_management like 'LOCAL'
          AND d.contents like 'TEMPORARY'
        ORDER By pct_used

  - name: ASM磁盘空间使用率
    sql:  SELECT
            group_number                             group_number
          , name                                     group_name
          , sector_size                              sector_size
          , block_size                               block_size
          , allocation_unit_size                     allocation_unit_size
          , state                                    state
          , type                                     type
          , database_compatibility
          , total_mb                                 total_mb
          , (total_mb - free_mb)                     used_mb
          , to_char((1- (free_mb / total_mb))*100, '990.99')||'%'   pct_used
          , free_mb/(decode(type,'HIGH',3,'NORMAL',2,1))                                avail_mb
        FROM
            v$asm_diskgroup
        ORDER BY
            name

  - name: 最近2天alert log的ORA-报错
    sql:  select
                to_char(ORIGINATING_TIMESTAMP,'yyyy-mm-dd hh24:mi:ss') originating_timestamp
                , MESSAGE_TEXT
        from V$DIAG_ALERT_EXT
        WHERE
                (MESSAGE_TEXT like '%ORA-%' or upper(MESSAGE_TEXT) like '%ERROR%')
                and ORIGINATING_TIMESTAMP > sysdate - 2
        ORDER BY originating_timestamp
