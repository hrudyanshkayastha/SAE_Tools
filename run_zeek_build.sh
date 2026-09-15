#!/bin/bash
cd /home/nysro/sae_zeek_build/build
make -j2 > /home/nysro/zeek_build_out.log 2>&1
echo "EXIT:$?" >> /home/nysro/zeek_build_out.log
