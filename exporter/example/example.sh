#!/bin/bash

cat ycsb_trace.log | ../bin/exporter -filter '[{"field": "aaa", "collector":"gauge"}, {"field":"Count", "collector":"summary"}]'