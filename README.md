# BloomSEC

This is the code base behind the Bachelor's thesis "Revisiting Authenticated Denial of Existence in Internet Naming Systems".

## Contributions

### Plugins

- _bloomsec_nsec5_ and _bloomfile_nsec5_: ZA and NS supporting BloomSEC with NSEC5 as backup
- _sign_nsec5_ and _file_nsec5_: ZA and NS supporting NSEC5
- _sign_nsec3_ and _file_nsec3_: ZA and NS supporting NSEC3
- _bloomsec_ and _bloomfile_: ZA and NS supporting BloomSEC with NSEC as backup

All implementations are prototype versions and originate from the existing plugins _sign_ and _file_. For more information view the thesis.

### Evaluation

All contributions are in the directory "evaluation" or "memusage".

- Corefiles instantiating different configurations
- shell script to measure message size, latency, throughput, memory/cpu footprint
- data from conducted experiments
- python scripts to create plots from data
