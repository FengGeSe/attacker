开始压测:  2019-03-22 19:14:02
+--------------+--------+--------+--------+--------+-------+---------+---------+--------+--------+--------+
| Task         | Rate   | Ratio  | Mean   | Max    | Total | Success | Failure | P50    | P95    | P99    |
+--------------+--------+--------+--------+--------+-------+---------+---------+--------+--------+--------+
| 流程压测     | 291.10 | 66.7%  | 3.44ms | 3.44ms | 3     | 2       | 1       | 3.44ms | 3.44ms | 3.44ms |
| 操作B        | -      | 66.7%  | 1.14ms | 1.14ms | 3     | 2       | 1       | 1.14ms | 1.14ms | 1.14ms |
| 操作C        | -      | 100.0% | 1.14ms | 1.14ms | 2     | 2       | 0       | 1.14ms | 1.14ms | 1.14ms |
| 操作A        | -      | 100.0% | 1.11ms | 1.15ms | 3     | 3       | 0       | 1.12ms | 1.15ms | 1.15ms |
+--------------+--------+--------+--------+--------+-------+---------+---------+--------+--------+--------+

