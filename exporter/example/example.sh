#!/bin/bash

cat ycsb_trace.log | ../bin/exporter -filter '[{"field": "aaa", "metric":"gauge"}, {"field":"Count", "metric":"counter"}]'