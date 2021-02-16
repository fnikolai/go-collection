#!/bin/bash

cat ycsb_trace.log | ../bin/exporter -filter '[
                                 {"field": "99th(us)", "metric":"latency", "collector":"gauge"},
                                 {"field":"Count", "collector":"summary"}]'