#!/bin/bash
numberOfCPUCores=$(/bin/grep -c 'model name' /proc/cpuinfo)

if [ length $numberOfCPUCores ]; then
  echo "cannot be found";
fi
