# Polymorph IFrame Generator
## Overview
See general overview here: [https://github.com/LimeChain/Polymorph-Contracts/](https://github.com/LimeChain/Polymorph-Contracts/)
## Badges Config
- To add or remove badges -> `badges-config.json`
- Note: the mapping works in such a way that a polymorph should satisfy all the traits for a particular badge in the config file `at the same time`, that is the whole row
 
## Genes interpretation
Each gene is 2 numbers read from right to left. Interpret follows:
- Gene 1 [0:2] - base character. Will not be morphable 
- Gene 2 [2:4] - background attribute
- Gene 3 [4:6] - pants attribute
- Gene 4 [6:8] - torso attribute
- Gene 5 [8:10] - shoes attribute
- Gene 6 [10:12] - face attribute
- Gene 7 [12:14] - head attribute
- Gene 8 [14:16] - right weapon attribute
- Gene 9 [16:18] - left weapon attribute

## Genes and their variations
```
const GENES_COUNT = 9
const BACKGROUND_GENE_COUNT int = 12
const BASE_GENES_COUNT int = 11
const SHOES_GENES_COUNT int = 25
const PANTS_GENES_COUNT int = 33
const TORSO_GENES_COUNT int = 34
const EYEWEAR_GENES_COUNT int = 13
const HEAD_GENES_COUNT int = 32
const WEAPON_RIGHT_GENES_COUNT int = 32
const WEAPON_LEFT_GENES_COUNT int = 32
```

## GCloud function deploy
```bash
gcloud functions deploy rinkeby-iframe --entry-point TokenIframeMetadata --runtime go116 --trigger-http --allow-unauthenticated --update-env-vars VAR1=,VAR2...
```