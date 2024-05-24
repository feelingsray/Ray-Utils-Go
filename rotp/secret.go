package rotp

import (
	"fmt"
	"time"

	"github.com/xlzd/gotp"
)

var secret = []string{
	"RAYRAY22", "PYYPYY22", "WLWLWL22",
	"KO3K24MB", "2JZ6NFRY", "5LUBT4PD", "22XILCN3", "36PHEA5T", "DVLZQZ77", "WHRRJK6Y", "OXODU7Q5", "F3LG3W4H", "OASFUZ5H", "NWOY7QPM", "BXQOQIN4", "ZS2XXNQO", "2FAWQAQG", "N2FKCI4U", "5XRSTHJO", "53MBLQIF", "FE7XWOOF", "Z22LA7QT", "4OJ66RYW",
	"P3IVW2HP", "XOWJEE6S", "MU65PYA7", "MDKP2NDM", "QWKTYYZI", "ELASR7YH", "POCIKIX7", "UFJHCYWX", "XLA734GP", "6X5SHRJ4", "ARQFRR4L", "M7XEKX4C", "INFIH5RF", "PBD2UYUS", "XN6PNPLJ", "HBMD24PN", "LF52UBQF", "WZFYMIP4", "CNT4MWUA", "OVZEE3VX",
	"OCWWP4GF", "6WEZN6NA", "ZDLZFF4A", "OY676V3Y", "CM474YC2", "WKSWUAAS", "RA2VNEAL", "BSHZLSHO", "VWK7E2GG", "RPR646G5", "BB6RIDXC", "2BEICIWC", "OTOH3IV2", "6B4KY2M6", "YV4RRBLV", "N5SAJC4N", "VJPTUXNS", "Y7Y2NKMS", "RBRZFTLK", "BVO4DRCN",
	"M5KD2IPY", "MRHSVNC5", "IGPZPSBU", "VCMMS44S", "XN6R7GE7", "YBYPYQTW", "FJXHQSSW", "IN6GKYBO", "QE3YUYFT", "Z6H7VL5D", "EFNP6O53", "IY5WX73T", "RQBISYMA", "PKGYDPGR", "AIM742FJ", "Y7JNGGYO", "KMCU7PWF", "SB7HMQHT", "OSMY2KUM", "HGW7TUTD",
	"AGE6MXS4", "2JAU5WLM", "T5OE2YRP", "UEW3THRG", "6XF2MQP7", "OHJQYAAE", "2BYE3PLB", "YWRTUQMB", "25Q2NXLZ", "CRUMYV56", "V54BTOD5", "KBEBM7BV", "TVKYEEBN", "CJXKCFIQ", "NRK22YXI", "5VSAU5II", "5SKRDPPQ", "VGXCPRS5", "JUBKIUAV", "UB7IBWPN",
	"M47NVRQB", "BSZTOEPA", "BQIHBB26", "ESAW2TYW", "NKG5C5IO", "57UAA27R", "EYBAHTM3", "NO6WSX5B", "2FGWDOVR", "73DI7GIW", "EO5PYFIX", "V3PUFZAC", "P5NT5472", "UTHKWF6S", "URW5ZMJQ", "XW44TVIH", "HJKDQDPK", "T5SDJFOD", "7TWKYFGT", "WC2JEV2Y",
	"7N24YAUW", "X3WPCWFC", "XQ3FPWXI", "CUWUUOTS", "4Q7YLFC5", "7GSOEIRV", "PCSDYEAR", "7QXWUCRW", "B4ZLOWCP", "GZZIOHAS", "WNQRMG2L", "IZ6UU6YW", "JXEHXGTT", "S5QJN7JP", "N56QSXWZ", "AB2R23H2", "7OEUBUUF", "BS6223DF", "J363O4TZ", "YZHOFYRD",
	"QGHS2MOX", "BUV2X6V2", "4WBN6VSF", "PKQN4XZI", "FWXGEW7H", "QLSVCHUD", "QYOZJDRW", "R4X7BHSN", "BUNVTQA5", "3EYEMX7U", "3YD4JBGY", "2DPTL2CT", "Q4KONUWA", "PZFRELSL", "XYMHK3RP", "L5H3RTOB", "YSPRJWNZ", "DFUGMGZA", "X2LEFXJX", "FWZ4IWEV",
	"6MCPLA6T", "22COER7K", "VWLGVAW3", "FKJWSS46", "CYDMLR4W", "MMTM3UUH", "VAHCXBIL", "2V7XAH4L", "4ZFWZKNC", "5UEKMWXA", "CJMAFZXY", "N5RUXGTG", "VTITQNS6", "6VB5X3HQ", "Q33LQDVI", "EN47MJ2W", "BTSHEMKW", "YRBWXZUU", "UH34QXDM", "UZZ4OJKH",
	"NVYQRRFE", "HK6XKZEE", "XSO2H22A", "QFUZBMZA", "PCYQR5SI", "6W5D43UU", "YEMTIQID", "QW2GL2TA", "KMVOE7SA", "4BNNT2LI", "4U5DLBJJ", "UK3DKTAE", "HIVLCUP4", "7LUPYQ6X", "H6CCWCF2", "HVKRODUS", "GJ2V6E5D", "OVPLLS7H", "Z3DL2TXZ", "GBZBTUXQ",
	"6D2SJQWE", "TREJCBF4", "WLEIRSNM", "GBJ76P7Z", "JQN65BXC", "ZGRRK63P", "36YH3KJ6", "R2W36XD4", "6SBRW2CT", "JGOFUMZX", "JYXUGGDH", "VHZ37KR7", "SSL4LZKM", "OOST6EEJ", "JEQSHJTB", "ZW62FO2E", "2MKZUFTU", "RASANQSN", "AUAEK6ZQ", "BDQD42SA",
	"5EQMQVRU", "YD4LBNZE", "RVCB2YH4", "VCQVYDOA", "LH7JLB25", "INZYERIV", "EBVLSU7Y", "A7W6E2KV", "KR6W5GKN", "3YNJYI44", "34U5MTHZ", "OQW3FYHY", "QEEASB5U", "TQ4GRUVE", "TIJ26QYJ", "I2SOYFLA", "2QR4LPW5", "K63DEZFV", "DSJDBGMY", "COJKTBFJ",
	"TYS2MVZW", "SETNA6EV", "UF3FZFDM", "XZJVGXIP", "2K5KBVOP", "5BJBQFH7", "VFGQDPS4", "SHOY57BU", "LNPKB64S", "KLHKR5FC", "2FG5EEPA", "MTNU5NPX", "EV2III4X", "AJHLGKSZ", "CEQUGJMT", "MYGS7M4L", "GFCK6HE3", "DRB5BRPY", "JUYD4XGY", "UYQFRQDC",
	"BREQFE6B", "4CRQX3WJ", "F3OGPGVK", "PWW2COQG", "PJZZSHYX", "4RQMVXTN", "7VOTQXDN", "XI4XNMJH", "7LEKQREF", "4LXK7C5V", "U7NA4GOD", "PZ5QL6HL", "7XDDP4SI", "XVVXRJMG", "P3M5LXM7", "GNKAI2SC", "P2JU4BN7", "7WWH6QI4", "B5OGXRHU", "J7OLONHQ",
	"PWB7XRL6", "6SYR3ZF4", "BUWZT4ET", "HAGJODKT", "X6V5RKWR", "BA5EAYOB", "O66YUDJ7", "IMGGNUYW", "YBU24W7Z", "IECNPAKX", "YZQRNOB2", "RVHFP23X", "WTIL7MTA", "MNR7SUN5", "JTASW2Y3", "JH7CZFTZ", "56FYROUR", "X7EJIRSM", "I5T53RNK", "AU2T64II",
	"YNNII4NV", "VK52LIGT", "ISSJDMGL", "MGQ5B2MO", "EBPRFCWL", "AMGEYABJ", "QKFL3K4H", "4ARKMGVX", "4C667Q7V", "6IGEIV7N", "D2KFFDXY", "P6Z5HPCX", "77YQL65U", "3SGEIHTX", "6F6T235A", "SBQMHEEM", "RP54GHO4", "SJMPCDQT", "CF3DFE3Q", "S3JXCOCL",
	"3XI2XJAH", "ZTVCIFJY", "5HCCFWAT", "KZ3JVRJD", "KNLNTP5G", "LNXAGZIE", "2CHRULPH", "24SIEGIP", "5PN3DIPT", "LUSLTGHD", "KIA6RUMG", "MGXRE6XE", "I2FZS467", "ZOFZRAWP", "5CC56J5T", "RWFFMPV7", "J3GHT6XJ", "AWE2WMDH", "FTBCHH3X", "MX4FMDIC",
	"GR6ZOOVR", "ERKZAE6Z", "VVZMDPZX", "7BCTTCRH", "ONPWZ26S", "MKXWIRXC", "DPCZQIUM", "GCHA7FN5", "6B6TCWI3", "446EDPAL", "3KK22LNW", "GGWJJCGH", "RCA5R2DR", "OHMUAN3B", "L2LITXW7", "OCDHEOQH", "BG7KZG5Z", "HY2RKHUD", "OIUQQ6SN", "Z5JYB2J6",
	"WBGYJ6MG", "6C73OXJR", "JI3TAYBB", "HOWCEPPM", "BBH2WGH4", "67NJZRC2", "PD54M55X", "IJ3PP3XN", "F3JURPF4", "RBEWXIBG", "YUYMB5IG", "JXXUZRYH", "W3GH5ZTE", "VBBKERQP", "JDKNYJ52", "ICG57H4E", "YWXVODW5", "3CDIF2DJ", "36C4YBOG", "TCN7NZ3R",
	"LIN7DRJM", "MMLWX7EK", "NAE2GW53", "AMB55P2F", "BJPQQZVD", "T6HAJLZR", "35DXRCW3", "3BSGUNRZ", "S5QOZFPE", "TL44MQKB", "HVWQISMX", "LZBXNJLC", "2B7LEBIN", "CEN6WLDK", "OKJA5D7V", "WLKCUICH", "SKJVYTNF", "KUSZ5KLQ", "HZO7ECJ2", "6W4PHODX",
	"3KKDJZA6", "2WJWMBK4", "SSVKS3YG", "62ORZTVR", "VGLRAIV2", "LKZUGAUF", "MGXIIHPC", "DUDLPCMN", "KXM5EZZY", "F3IS56KY", "47GFUWWE", "54VYGCBB", "EDR33ZPM", "4J26DQ37", "V3XGKZO7", "M7UGRR2K", "UD5IWUYU", "UDJANMU7", "MOHPALP5", "ZDNDS74D",
	"5V3GQESG", "XSS2DMOF", "VWOBZDKP", "XM6REG7F", "QSNFHQKD", "PUJXNHXN", "FR27T7DY", "4YGCJ3AD", "HNXWUMFZ", "PQYJXUAW", "2UB45M5B", "EQ77DH3M", "O5TSHOVJ", "XCFGQOYX", "KRRKHGWC", "HEP42SSA", "WT4QRJOK", "J5YSWAMV", "SDRGCW7L", "656ZF5JI",
	"PRM5DSQM", "GJLQNBBR", "GXPD2VUW", "J7VYTTKW", "ZN6L2NHA", "RK574VC6", "ZL5CRRRZ", "2LWWE43X", "BNQZ3UIB", "VEDMFW5Q", "WGSAI4IN", "JO3TPYVY", "OQWHJKXO", "AOWW4QTL", "YXWMTJH5", "Q4W7J5VS", "B73VEE3Q", "5FGYJ3Y4", "EWXMEN5J", "Z5DDJFKT",
	"IBATM3XD", "TFCYUEJE", "NBA334WO", "43XPUPZE", "4ZGGIZTC", "T6BF7RRM", "HED4IVV2", "UN64ANSF", "L3MPEAPY", "37OXSUH5", "2NZLZMUP", "6XCO5ERZ", "T56CZUWI", "TY5V46RF", "KAJZBXOQ", "P22MNYTH", "DOINNY47", "2KCUSQ2L", "LAVEOMNA", "LDMXBMX6",
	"2XNMUCGF", "54L7Y2SQ", "6BYS3Z6N", "KDYWS64B", "CB6KFEX6", "NNMNK4UJ", "CT3BFPY7", "BXBFI5FG", "EULI5VBQ", "KNWM6Q3B", "CJUSB2W7", "VGGEPO6L", "VC53SZZJ", "TXJ7RPRZ", "M4IOU25W", "EJHCQPPE", "RGOUD22B", "4CSNUUSS", "56J3XAOQ", "BJZQBQS6",
	"EUXYJFTH", "ESGL4F5E", "XIV7YGCT", "PL4S3Q5Q", "X62WY6ST", "7DWJ5W76", "QHN6YMF5", "4KMR4RQ3", "K6RJ5RC3", "LRQI4H2E", "XND55YY3", "P7RF3CA7", "H5REMUYH", "OKPNFEJQ", "JMPB6X56", "VQWQTAI4", "7M2XDERM", "OUGXZ3OX", "OIR2XFV2", "GKKNKQ7Y",
	"SO5VSYRZ", "WCLZPGX3", "OGBMSCC2", "QWGMC7KJ", "W2CPJYIU", "Q4BH55HJ", "6DLKETDT", "ZHXKUGMD", "F4G5XQXB", "ZRUBUA6E", "BTKUXLIB", "BRL4JB7S", "BPUP4M2P", "QDPS2NQL", "B5RG5Y3I", "AEJFNWVZ", "TKHJTPRD", "IRZR7A2Q", "YNVFCEVO", "TUSEUW55",
	"LISIJ44R", "LOR33GWP", "JQJDNVQA", "AUFGSN4K", "TII2EDZR", "LEJOYEUP", "KMBNJA57", "QARRGQEC", "ARAWFL63", "LEOJCYF7", "CL657JHM", "F54B43OQ", "3362ILX4", "GP45GJ5A", "2ZMRQRQO", "TP2W7OWR", "VDT6TBIP", "3DKTO76O", "VIH3WTRP", "24T3UVVS",
	"JQCDSRAM", "VY52X3E2", "L4S67CG3", "MOVQNLCF", "N2QYFUDG", "SINIC62J", "EYQNPDSV", "5FFVXZUW", "IWJW5OE4", "U5JLSSFP", "QVKADCTH", "P3BOS4Y2", "X7CDFEVA", "3USGDS3E", "74O26GYC", "L7ERBNTA", "2UHGEYAG", "5IUFBEFJ", "WOD4EOPH", "BKCBMXCI",
	"BUMWF2ZQ", "6WVHWFH7", "BK2LID7H", "6QDDLKMW", "ODRDHYCR", "BNJIDPJQ", "RLJMWTHE", "B7WUER6H", "ZBWJ2O53", "NL2KYV2F", "Z5XGZ727", "JWM3L3KW", "JSBQX4SD", "ICESXF6F", "IYGH2TLL", "5J72TGQC", "TAOOGMKA", "RAGDYIYP", "IG7XT3NF", "YMQOXKYC",
	"AGSDZNEJ", "2CZWM5PH", "YXPWUBAI", "L4HKPFU6", "DAVBDO74", "3VDBBPG7", "3ITFOZN2", "SLIYBEGX", "OP6AJMIB", "SDMUH2P4", "2JCHKJKZ", "2375ICGE", "7AWNDDKS", "3FDEJ2X5", "KW2FJZCW", "S22J6TAS", "5GQQTIGN", "OIPH72L3", "NUMHWRYF", "NI22T75J",
	"TUFSMIPS", "VYKXVGWJ", "ZH2XSY2M", "53Y2Q2QQ", "GH7OTE3O", "FYPVE3U5", "QFMW4QX6", "QPQ7KPTB", "2DGUCUTK", "MHATK5WL", "7MUMF2ZB", "Y2L7P46P", "5TNTLQCF", "AXLTD2DG", "UF2I54H5", "7PL3JMKT", "MSYT2EXZ", "4YHHHXLH", "X2N4IBIV", "6GI4QW2X",
	"XQ2G52TD", "O4XFGNTD", "X2X22CSA", "JGOOCXUA", "6RFFOIJX", "XXWZXK4F", "6RJOT6A2", "WDUCWGNB", "GZ4DUGX3", "WNKLEACV", "OP4N42WG", "WXJTSNLB", "GRY4SHG3", "WAFMBGRV", "NI7WYLEH", "MQDL3W3H", "RC6AUVSG", "IMQUAVFU", "RC656RRN", "LA6AXU2H",
	"QP7OTCMM", "TT5O3MON", "7XJVRADI", "Y72OIGI3", "ITLPZCBU", "TFTTZ5MN", "A34M4UJU", "SVO5IXBB", "CHFWGWL2", "T4XHTDUG", "3EX3PGJG", "DI3VIEQE", "S4KZID26", "FT6MZGFV", "FNCWYGDY", "AA5ILIXG", "GJT6GOLF", "VIAHVNW7", "ZPN6WVXX", "MW3W6FE2",
	"Q6F5WH4C", "EIBQXV3C", "ONVF4WN7", "3GA27F2G", "CJFORRGV", "PHFTGFHJ", "L3GGJRDP", "25U4L6QW", "U5LW67TE", "63L3V3SX", "6NR7XTAG", "OHPQWOLA", "62VVGNVZ", "NQD5EIQS", "2GAUFUBM", "WLOIH562", "FAQ5KMKB", "72S54XWI", "BG7GJGPT", "QADH2B2N",
	"R2HRJNHP", "RG3SWQA4", "PZ5HYHMK", "ICPQ7JQV", "3X7R6ILO", "JLGW6DVI", "INPNARSW", "XHYS7O4Q", "5REI3TSH", "EFJY5WOW", "RJLM7H35", "DV5WMRDI", "3OVYSMJD", "LQX4UQWK", "XEURX2UR", "WKYFZLBY", "GDGOZFLR", "SWID4SXA", "35IHQVWU", "LM3I5E6A",
	"BW5SNE4D", "23EXPQYJ", "FOFLRWFY", "CRS7E7B6", "DCOFXENV", "VUF2RYSK", "WPUKP73Y", "YEYO3DD4", "UI4D72BK", "E4GHSGNR", "MWL3VS2Y", "5I55B7CF", "EMDW2WJD", "UPQ23QT5", "4LR4CBHP", "UUEP3EK5", "HLWIHYOT", "HVN4B3RB", "4L6PNLWX", "P5DE7XS6",
	"UR7ZY6I5", "LVWZQCK6", "DQVDITIK", "JJEK6G5E", "GON7AUKL", "6MNDXFLH", "OE4MVFVB", "Z2CAXPBH", "AEJWH2QB", "YXT6PPSU", "43F7462B", "APMUWQQ7", "N3TPZR6G", "Q73FCXFG", "VGPPZI3B", "UM4SAEHM", "IKYCF4UX", "XOBF4USB", "IEKMVWRZ", "XIEO2O7E",
	"7SQK3DJ6", "4YCJUOKV", "4MH7O34D", "G5HPHE33", "GQLBTEN7", "67GLHZI6", "ZKNSA4HV", "Z7ZFM623", "GNDEFFZ3", "GCALQ2K7", "7HYJI3LY", "Z7JT4CWW", "NBZ2VLVO", "FX6NCXIT", "YMAM2AHS", "BZVSFAYY", "BCSS6EZP", "ZWU7GKV3", "4OJFB4CR", "JU4E2MAI",
	"F2SYOW3G", "RA4CZV6M", "FZ4UX33L", "YM2HUQSO", "ASTGZI7Z", "4YBKQ6MD", "T34NVWIO", "3JYACI23", "LXDKACCY", "L5LRYTRQ", "TPSAS6QH", "KE7XP4XK", "KURHINWD", "KW2ZNFDN", "PEV2TXWT", "OIFNWASR", "DM6A3ZP3", "SQJDBRMG", "KUHCJMIQ", "CESV6EV3",
	"5QXT7H7E", "ZVBCZS6E", "2IRKRYO4", "4IH5UCJ2", "6NGMOSHR", "KVQUFVSD", "M3QL66CD", "VNGJW5R2", "X3QQPCQS", "UZX7INPK", "UPERULCX", "UHAZFYCQ", "JMHP65CH", "U7BOXQSH", "FP7EPZQA", "ATHDJ4QY", "QCE2TRD4", "D4I3B4LJ", "LQMOLR4W", "NPMVEE4N",
	"ICCU6N4F", "MIMLWQM5", "HSBA3K3C", "AQAPU342", "PGECPP67", "K2WJYSM7", "XWNYL7I4", "64W7SSEH", "PF7RNE64", "7D4IRMZ2", "3HKIWEYF", "TLI7ZTRC", "6ZO6SWC2", "GXNRF6NY", "62OK4T7K", "6WN5Q6KI", "S3YQHVXS", "5GWQLNG5", "O2LGJ4QS", "66HNOY44",
	"G7IOV5BP", "7EHCYJMN", "BGZGFTTY", "PWUOMLRE", "C2CBRDOO", "5NQKGYBU", "Q7VTVXLO", "ZFBWMPJB", "VXMW3KBJ", "AFW4CCPU", "ILR7HZ4H", "YFNS477N", "3CWVBXMY", "JMS5HOLK", "QS7Q6GYV", "4A3SD5E7", "QEAR4FEY", "DOZGXD2X", "ULFQWCVQ", "TOBO52R2",
	"LEPF7EMY", "TYNFE4JD", "CADLUPMZ", "D7LSDKUK", "35FMLKTV", "SDRAAFAA", "65C7I55K", "CFDXACW4", "NLNJFFFH", "ZJIQM5CR", "5TWUBTP4", "LZSWIKMH", "A4SH4I3D", "M4A27TVA", "46Z5HKSL", "SEFALCPV", "JYDJTHSW", "CKXPKYHS", "NWVSRQF4", "ECBBGHBH",
	"N6YKNBFZ", "EBJPDSUN", "F5BBGAPK", "4DKFMUL5", "HJFUSMJI", "BTK2WV3F", "MRFBMLYP", "ETTUTDV2", "7XPXY3BE", "43QLKLOT", "T36LSDM6", "P7Z6X3JJ", "UUBVP2ZA", "OHCJBFGH", "L34H2PF7", "OZFLBFCS", "O6R6H4O4", "3IHSYU5D", "7OQ2AM2N", "RGAY5FAO",
	"OB325I3Q", "ADZPXHSQ", "NTJVQMAI", "B5UZGHOT", "4ZTZ5WA3", "S54AUOOG", "NFXTZG2R", "IXBCSPZI", "M5B2JDJE", "PRP6GY67", "JWQ565PY", "QELAUR4K", "NXVUW4ZR", "ZZC35TW3", "A66PCLSG", "ZU55KUB6", "YIHSNKOF", "LUDFUDNP", "L2QIZ2K2", "L7O7SDZ2",
	"YELDGNGA", "HIVDLFDL", "3LBGS5PW", "HRB2GOQK", "TP6B5K5U", "LYJECE2H", "V5DHJ4HR", "7JBLRNYT", "CPAKKWZK", "ATLN6NVV", "2XWHDPJD", "HRNXGZUA", "2BJ6MRBL", "WDL7U3SV", "KHFVM6RM", "S3NUFHCE", "KI4IRHE3", "KAZ73VX7", "VDB6UAWX", "SGXSXJT6",
	"MO5KI6GI", "ZUT7KNCW", "EY4WPEPB", "26GFIW7Z", "AQY5EKCH", "GNT4ICA2", "AD2CBL7R", "UZJG7ZVM", "Y2UZWQT7", "VDSIPTCX", "JNL4ZWWF", "MTHTQNDQ", "HFHCIUDQ", "ILC2SHH6", "F7BZL2GW", "XDJAEAVW", "4S4TAU2E", "FFSDYZJ3", "4M7W6QVO", "PWPJJV35",
	"4D4Q6IHH", "UEOHHP7P", "YMIKOC42", "XQT5U3ZE", "OUKN4D4N", "7KQUUU2F", "2MITN425", "PC3GJOOT", "DOB5CVNM", "PUOMHYKW", "6I4QEWRZ", "KPGTLO6E", "RJBKBG3P", "G6S2MOPF", "SD6BRH4Q", "GJXAKL3H", "N52YSR4I", "ZBYLJI3T", "NH62CXKT", "GRROMJ5B",
	"AXYBPUI6", "RBZZ3W4U", "6PPYT5LN", "RT372VIX", "YKLTTLNN", "L6ZGRZUJ", "BBEWK4DI", "JEYPPCTF", "IAESU2SQ", "I6ESFLYB", "LKSVKDGL", "HQR5AIWH", "ASY42PE7", "RISQTBJN", "4ODXMUJM", "QCJGYHM3", "TH45QKMT", "33XO6ABW", "AITQDY7A", "TTTF6WEA",
	"HBA47AVB", "LJIUEWDL", "D5XXBJKG", "E7AW3NYG", "QJLMEPQP", "OY4AARF5", "RMLE6P2A", "2AQKW23X", "LT7LTJ73", "TSGSNJRT", "HGTVK6FW", "KSBI7WDB", "CGPH5AJE", "EUMOVFI4", "2HGCBZNJ", "FNRGWRZU", "CVF6WLM6", "MI55PPM7", "EFKQS3W4", "NNSXLBXU",
	"FF3ZPLZB", "KJQB6CSK", "4HHAZLRK", "4LWH4WMH", "NA4GV5N7", "3T3ZYFY5", "FEFBRIXV", "ULBDM4DU", "TNKXRUA6", "ED6LNWEU", "PIFBFVTM", "44KUSENB", "X6F3KN6B", "PCTKNYJ7", "MT2BH5IW", "O2CTD3RM", "PPQJ37QE", "3EQ2Q3R7", "D2VRMIDE", "2CDXT6IA",
	"6D3ZIV63", "3BCQBZ4T", "WFEFEBJZ", "PBFYGLEX", "WVAZBRKW", "OXQQ2W2O", "HVKD5BEM", "CLZDWEVD", "KAPWJR7A", "XTPNCU7Z", "BYX4F4KX", "WXXC6DJW", "W36WSOEU", "IQ5E2XDM", "ZEDYN55J", "J2N7GGNB", "GYMTJRY6", "OMIBENKE", "YE4HKYR7", "BGBN75BD",
	"ONTFFRBY", "DLMRIUWP", "5PHT7MT2", "IDG3HQV3", "3LR2MHSF", "QVDSYKH4", "33IRQUWT", "KQWZ6U5X", "LTUIR7WU", "AB64L73C", "DL2TSWIV", "4SLD3J4D", "HWWKSBZO", "TUEOQPAR", "CWRBVGM4", "7KTVHUKC", "DOOVPNWN", "3EIMZP2D", "HGTMOGIO", "KKKUGLYO",
	"VOTIZWV5", "DUR3ANSI", "VBQPKQWW", "7NNSBIUB", "VT4GKDIX", "M7KJB4VB", "BBB5MOYY", "NJPQSGGD", "636EMKJQ", "2JFTFTIQ", "5TZHREN7", "U3U2G4KJ", "M3UC47YE", "JL5BUIK5", "IWMV64ML", "GEZYGUJV", "4GIL7X6M", "RLGPGO2X", "M6DWNTN7", "VGMWUJ2K",
	"HCL642KL", "4NEDEZIU", "CN5DKLSF", "WK5VGQPF", "JOGU4I4Q", "B244T6OR", "BXSMXHKP", "GDIQ7RKX", "RQCSGIXC", "IS4HA2MQ", "Q57NVC6U", "BRJUTBIP", "INEU2ZFZ", "AXF5RKZL", "R3DRFSCI", "IZEFZOS4", "HDE2EMIG", "FXPJMD6R", "4LMRKFVU", "KBMAPOPS",
	"5QAFFCUK", "IEPVDL3N", "AEYC6XUW", "BSAKO2M7", "FYLNESJJ", "USUIFWC3", "UKUXIA5Y", "7662B7ZX", "OPJFEF4F", "HYDH4ZRP", "ZJEJFKMB", "DXSM3CJN", "75NPBFFX", "W6NAWAET", "TS3T3IPQ", "PAVV7743", "Z3DKBLKC", "B3P5ZDHN", "67IA53EX", "XDIIUPDK",
	"6HU4ZGPV", "YCFPUJFM", "JKP6OTED", "WO5FGUS4", "IIF2QVHR", "5IYYJ6GJ", "Y46PCFFC", "EOMS7PMF", "BUTSYYL4", "EAPVPQIP", "U4YZ6PSJ", "USNRAVNG", "Q6LPHNLR", "YYW2O7QL", "QWNOB63J", "LI5M2P3C", "DGAQPL2V", "AJKVCTH4", "IFW4HJEG", "QPJCNSGL",
	"TDNB4CO4", "IRXVIGRJ", "LXE4AHCJ", "7RX5NZ2O", "KNVFV53W", "KESUP6M4", "3SJQAMPF", "KOQRAH27", "CSOUF6HK", "SYEHZKSH", "CILQZGNB", "GQFT6ALM", "WS6GBFFJ", "5R7OZZBM", "FQ6C4DMK", "CUIBVH3B", "NIUETKBE", "E553L5R4", "B3M2ECRU", "NQMCCDWX",
	"KEBB2XWP", "DHCFRQWD", "5B2MXBL6", "QND74YYJ", "FYIHJ5RV", "Q6R3QA6I", "NZ7OT7I6", "FFIAJVOY", "4HRTPLMD", "YE43FDYW", "4NS53W6J", "6ERQPAZG", "XFJ3VXOC", "YA34BIHN", "E4YPVNSL", "3CFC4E7W", "XOA2DJP7", "Y3S24VQX", "7BXMHRB4", "VTGETDWS",
	"OZUXWNBP", "H3VYKJAE", "GDSKRCMO", "Z5MK24OQ", "7AH7TFVP", "FGF62YR2", "B2DFT6RZ", "GVBF3TT2", "JNH4TZSS", "VR23M7BK", "BGX6KNIN", "O2PFCQHG", "BYVU3ZHF", "VKGIWN2T", "N2D7CMMY", "RG7DZP6B", "OJGTM33H", "ZOPKTWIT", "AJKNXOV5", "3MWCKWDM",
	"4JWV5B6B", "IBXNNXQL", "D7NMGDBL", "JTFBBBWC", "I5YYJHPL", "ZQXLMPZI", "3LI4FOPE", "LEYNUNZ5", "2Y7SUITX", "X6XZMPEP", "Q4EHAX7M", "2KL7YGJG", "C4ZGOQOB", "GKEZEI3M", "FUCBJAYW", "CUWP5KTT", "52VY3XFJ", "LAR3QLB3", "5ERCJQBT", "2RXQBX7L",
	"656W4VGK", "WZVJPBQH", "VFQMGX5S", "5HKWOOCE", "C5SEGVB4", "E3IL73AU", "BQKKX5RM", "WRKTMX7A", "Q2QRFKO7", "UMII5T6X", "4SLHWW6P", "XQKLNKMD", "MCQSVU4D", "VILBOE43", "I6KIHL4S", "PRSEB2BS", "I5DYUG5Q", "X35235Z3", "NBIZPUHF", "4T2CAQD7",
	"3ZCBRO2P", "KFKABRUY", "ZPWD5GA7", "KDE7ALFW", "GXHTCYD5", "62S3L23F", "RGNPFYRE", "ZSJSKQOP", "DXXVBHLZ", "64SEGAYE", "JQTZJOEL", "B64QC7ED", "YQOP3CEC", "ZHSFHWWI", "ZCSG4QU3", "JRQMUXUT", "U4Y3OIUT", "JBTCHQEL", "GFUXK4PS", "F5PN45D7",
	"ODMB7IP5", "I3JJV6FB", "AZHYY2Q7", "ADQ77S5J", "EJMCUY2U", "YMWFDZ2I", "AQA5INHT", "3SH77ENO", "QWRGV4KZ", "XKIXO2QX", "2HR6ALJA", "AB56FCVK", "3ZME5TUK", "T4UJTPT6", "I6WIMQUW", "PM32YJFD", "SAJRRME3", "L7JGFSFP", "KZAZI4PM", "3C34EPLK",
	"7Z5EDFOV", "SF2NMKP5", "FTKAILUM", "NV6UR5IC", "34GDD77S", "UPVXGKKQ", "2WB6FETA", "SSQNIMN5", "ANXBLKJ3", "4KEV6VDR", "NAVI3YIH", "EUDMYFNK", "RW27LDIH", "ZLZWOLCF", "YKZJRS6C", "UNSJDRWL", "5SPRHJD6", "V4OATJIM", "LSIUN3MC", "473IZA7Q",
	"WQ4MLJLX", "IJP2MTF7", "M54NJBZ2", "D3DA5AUX", "EWSUAK7V", "5CTDDP2T", "K7KWWZUQ", "3UKOI46A", "KJECCMCX", "6MWVNPVE", "DMTRYUFV", "7QUVOQDJ", "VSLJYUJ7", "CILNGUN3", "JW2RSUSR", "6IIYPBYM", "P6ZJJUMC", "VLWRBJ6K", "6LXF42AZ", "RZRYG4WP",
	"JDR55XEC", "G2SAMT4T", "YXJQP5XQ", "QBICC4CO", "SQU3TS46", "AVBNW5W4", "AMEQW6S7", "J7JJZFE4", "YTXNXH37", "QKFWMHUU", "3AXKXYIC", "NUF5FONG", "ZITBCXUI", "42OASKNR", "7Q4IQYTU", "MGI4NGIY", "YJP5S33V", "I54BQNBX", "2FMNQ4UB", "DFMREDTW",
	"PJPGXLP4", "7F6GHK25", "GLS3D4OL", "2VLOM5RB", "KLZCKKY5", "TQZEOXE4", "ISX4R5NC", "57LS5JVO", "7AKJHFK6", "4ALN5B2R", "WIDNM4CD", "K4RUK5YG", "LQVT3VBV", "XETYZLHS", "P6IGOYYG", "7UMI7TCH", "CHJPOFKQ", "7LKDC4I6", "WDVJESLF", "ZROOWREU",
	"KLB6PRWL", "X7UWFLJI", "R5LKIZTF", "PHINNRBQ", "RCUP5KM3", "MVJX23DV", "BCZLEQWM", "YKGCPMUK", "LODCIQUT", "BUUWDTYB", "Q2TJXRT7", "4WLBFIMP", "4SMUJTXM", "5UET2WQ6", "TVDLNAK3", "BWQKEXHF", "MNYRUPAW", "KJFRUBXG", "A5WJE3AW", "YJ64HD3N",
	"IDFGDTYU", "I5E22FZH", "2LYOTMLW", "COKFNLDG", "XI4ZH5IU", "YOKMT6MK", "VAOWA6ZM", "2SPKNT3C", "6G3J4OUL", "K4DCNJ53", "K6YB6KVL", "SFPVI42C", "SX23QQPN", "CLJMRP2H", "U72B3QNV", "4TIYIOTY", "YHV3GFK3", "FFEPZPVZ", "V4L74VPW", "N72SA6LU",
	"VDX2XGLV", "2G3ZJEUF", "PFS4MM6D", "EUDEMYSN", "ZGSYI2X3", "U47AXWNL", "MP6UDJSC", "PUXI4MHQ", "7KFL2O3T", "U4WAGORJ", "HTUTE4XE", "SXO7XVFL", "X52FNDZG", "QDTVXF54", "MJBJ2QI2", "WLD3NOTY", "XB3D6J4A", "5TUXYK7W", "WH2H32N5", "6AYRQJWS",
	"WCSZY734", "ZG3JHRDM", "J22BF4YQ", "VWXUIHDN", "FUNDLV6L", "6SNW67ZJ", "BPZOQXBZ", "5KIBTB4O", "FXFQGCHM", "WT4EJLRJ", "23GKROJS", "J7OSB2JF", "5ZBJKNM4", "U7S5GOBK", "DTA5DPHN", "I5XVORK3", "J2MU7LSL", "7HFIYBXB", "M4SQYPMF", "DOC7CBBS",
	"UETFCA3M", "YNC43A7C", "CKRQPL3A", "CQXDSXG5", "3UGXG6Q3", "3KV2D7VW", "NUW7SOB5", "GTXEIBBY", "VHC4GWGU", "VZA7DVMX", "NNN7BAD2", "NL74ORWA", "5A55L643", "RNOQVQ7R", "VA4UD6GU", "FVKIAE4Q", "ZKK4MEQG", "QUFQGXE4", "UKUTEELY", "P67LBSA3",
	"EQ4CBAB3", "AEJF7SXW", "VX5Z3CKN", "7DDV6CYU", "QLSISPER", "XFPLHFR4", "ULE72V32", "ET4CIS3G", "SXNVL5WE", "M3HYSUSO", "GZVMV7OL", "GT2ECKGY", "XHRUWVAW", "6W23KM6A", "PS3K5XZ5", "P4NCKKBK", "IKNWOV3H", "WGCGB3HF", "UM7MYTUP", "GOSBRJKP",
	"WEAFOTQS", "CGNYVL54", "H4UNORS4", "BQLAR4NZ", "OVKUUDXX", "Z3SX47GC", "6NZMWQ3B", "GRG73UIL", "AXMOEZIE", "BCGV5EI3", "BCL5BURI", "5UPART6T", "NH7TE3JR", "DLGGIGEP", "IGCFQG7A", "QMJZLWCO", "YQG4QOPZ", "IWVPT2ZW", "MWA3JNB3", "LSA33H2L",
	"7F4RLJC4", "XILXD3XN", "QL4HWWBL", "OR3ZZCMJ", "KMSN4MIG", "KHKDNEWW", "DVTKNT6H", "SBANELLR", "277BXVWP", "LXHWING6", "D3OJMVB4", "GS6J3GYN", "6NFA7TTK", "THIIN2C5", "KODKSVPI", "LK26W62G", "OONU7PK7", "KW54QI6I", "VCJDW7LT", "YIRCOIKL",
	"NZDN5VJX", "WU2BA3EU", "ETXUUG7S", "EZODWOKP", "FAXNT7UF", "NFAO2DIX", "XHWB5NDU", "PRBXFP34", "N6X5X7IU", "XJORKHES", "FHKUB7B4", "IKXXGZ5H", "HTKQL7AE", "MZXTRWNP", "H5GSKXMP", "MGPJED4G", "POD65BBG", "UUP6OAKO", "OGXF5ZS6", "LMGZAZN4",
	"BH73WPEX", "CD5PKT6U", "OGMCN5JS", "PBQLNAUV", "ZPY34XNF", "3HRWREST", "KBDK2XFJ", "ZT55PYIP", "BYKUDHCM", "QVITX67X", "JC6KKI3U", "BUSMPCMZ", "6SA7SNYX", "N63S2EEC", "V4LG4LP7", "N43OECTS", "CRKTCDYN", "BYGSR3R5", "I47ZYZOI", "IGL3FWMT",
	"LIXSXQE2", "QHTC4HCG", "A3YYNA75", "ZHCPQLK2", "DJCNKOKT", "2Z5URKX5", "2SVQOAZ2", "74OUX3AP", "GK4UOVPL", "UNC6FA64", "UF4TYMMU", "V5EDHIUE", "57V2KSPC", "4W56MAIL", "7XC644B2", "YYYV6GMY", "EMXISVHW", "EOG4V5BT", "NTMLYI5R", "5Y26LOXP",
	"U4ZWUTIX", "VOU4TUW2", "7EJLW3AX", "PJI72FLW", "LFXSMRGS", "XBZI65EC", "XJU6JZAJ", "4LFCYMZS", "7PAF7CG4", "XVEOGUKO", "GZ7O6MXZ", "Y7FQWUWR", "SRFXPZXQ", "XVVMIWNH", "72AMPS22", "O4XHC44A", "ZK4WTVEQ", "ND7ARWTT", "FSQVTYZ3", "WI5RGD6S",
	"WP5GPS52", "CNLJXW2M", "5TGQMOIY", "6HOPGSGP", "EB4GOYQB", "FDIO7UJR", "DBH5C6DO", "45WQF3PM", "4PZWDF4P", "YVGNWOXM", "YYNBZYSK", "AUIYCNVL", "B6D6CWAV", "24CCWR7K", "3OJJOR7C", "TKBLFMG5", "22GWBXPS", "IQTVSWIC", "WUBJF5TA", "TZB4IHO5",
	"CNAS2XNO", "KRQJ57XK", "KDO4QKCI", "WH6RCSQA", "NLYFYKNL", "SR2USROD", "B5ELXJLO", "ZYV46767", "F4NCPH5P", "453EGA22", "56BY2KUX", "MGMUEFC5", "KKWWJ57I", "APXKMI2G", "MEW6AIVE", "Z7IXZXKX", "4JTJOOXD", "XJLAGNYC", "53XANGFN", "MRGI4EQH",
	"MSA3T45R", "7QL3YUZ3", "6QQG5HBA", "HHTRSL66", "7NCEVWI4", "QAITYUT2", "DFC5VDGA", "ZSHQYJQ4", "HZGXRSQV", "6TG3FVPI", "GCSKMO3T", "7EMBEV4L", "OIQGHBZ2", "BQ3J6UGE", "BUD4RDRC", "WXVB6TJO", "J5QTTLFY", "WROCLKUQ", "BEZXOWBX", "WY37H7CP",
	"WRQVIXCQ", "YPOUA3TI", "ZF4QUFI7", "7Z4FNDO6", "IX4V6VFO", "ITJLDYCZ", "JN2RQ4LF", "RMAETHVD", "TAQXYNAB", "ETUMU5ZN", "3V6LKUGX", "I4EC56QV", "MOYDKSZB", "HSCXAKWM", "IYUN2JWE", "EMPNSSHN", "3C6AV4SD", "Y7HUXJNA", "YVVQJKSX", "3KWAVMXN",
	"DFIJYEHS", "DAX44PCQ", "BN6MP35N", "TPSF4MFS", "LVRY7ZQP", "LTAHS7KN", "S5M44WRM", "WMYEENOW", "W7WTXYYU", "5ROIAW7T", "VOV3TB2R", "2BTCMFJI", "2ATHCBZE", "5E6VX2XP", "3SXMPJVH", "GG4BSZSO", "NPYBNWZM", "VJFO7Q6D", "FAX7MAWQ", "HP4W52OB",
	"AS5YUQBL", "Y72LH24I", "D3I65R2T", "NHTTXJPS", "MJIKI7HC", "EHO53KDA", "74QCNW7G", "N7NSFKRP", "ED4FJJMM", "G2INK5OF", "7SEEZAFW", "H67ERBH7", "UOSVXUZ5", "3QGKQ2Q4", "UZNZTDKZ", "DTU6E7US", "DO3VHL7Q", "XHC642I5", "GLQCASXJ", "6PMBK2IR",
	"ODZGZZCL", "OVL5RIBD", "2BG4ZNED", "OPPESTDD", "JJQIYK2I", "J4DNXIIK", "GDYFZDYL", "EGH44NTJ", "BMNQ6VPP", "K6TEBD2N", "GERYVX3B", "RI7X4PIL", "Y5DMVVOK", "M7OP5JLV", "FCMXESL6", "UML7S7LK", "TTXCYRSO", "Y3AJOJ7Z", "ILAK6IJS", "WROBBTEP",
	"YUKV46RW", "IN323Z4P", "CSZQTCLH", "LQZEJG23", "PETTDKL3", "3I64LQ44", "6NZ7QKIH", "SAKS3M55", "5DURUTMU", "DYH24UMV", "JWRIV5NW", "DCLQ5SOW", "5AKDAYJU", "N27U4IBA", "OX5LAH46", "CCQNMYFK", "FM3QTQRU", "BXOW4QJ5", "F3I6DKXI", "JZAHNQZN",
	"FWCY7MZ5", "UO36RYXV", "YXFHPBIR", "P2E2SCTP", "AS4QENPC", "PDVGQXDW", "RGUJEL2K", "JJV22QI6", "ZE2PL5HW", "WSHX2W7G", "4MXYYGXL", "JKXP2EDI", "EGGQ4D5B", "MJ7UCUBU", "JL6ZWOAP", "BOWAPAAH", "JTUQYJRI", "QXBH7B6S", "72HY54JM", "QRPMQHTK",
	"YZ4VXZZF", "2ITJDMN2", "RK75HCLG", "L32LBEWX", "HW7RB3S2", "TS6FVXSO", "L6KEMO7B", "C2IMEDQC", "Q6RQZ3NM", "OUFDEPBC", "HA3C6AA3", "6APODTY7", "K5WBG6D4", "2PEWYZIT", "VRUNJSBE", "FHFOYO35", "GDCB32V2", "SXOV6GCB", "MBPMCQ57", "2U3RALIY",
	"4IPW54AF", "N3PKBJLC", "5Z6LBIX4", "UDHNGAUG", "EYMSJGRN", "B6MGMS3L", "VCLOCG2G", "X4C4FRUE", "ERDSHBBK", "55CF2H4H", "HR2PB2SD", "7VZDUJ4A", "5ZDWHR2I", "NNRGZ7UF", "GPS7OQFY", "NTNCFPRD", "PZ2WGW6K", "L3KK2BJH", "SQL5MOGW", "6TSN7YAT",
	"SBZSAUKN", "B3J36QVH", "72SKQJ5W", "OTZT7IIR", "6GAGBQDO", "JMF2E2BV", "RJO6LTO7", "ZSPS7QNT", "HSKRWIK6", "6MLFAI7U", "6KGMZJ6M", "GGWQF3WY", "FZDQX2OI", "M5FEZFLP", "OOVNZBFI", "4NUA4LQG", "SZHWVJGF", "5JHLH2UV", "IHQRMJ7K", "YV7VNIKE",
	"M2Y7S5JI", "QHTHQ3NK", "3A7LWRHP", "ICISHLN7", "ZYYFKTI5", "CRF4KQJ5", "LPKERBKS", "XHXOYKQN", "3YTDQJWM", "74GCJOVF", "3KRFQJSP", "D7RNFMSD", "2H2BPA6J", "7PCLGPU4", "ZJH53JVB", "2IPVLJ5R", "SG2IBC25", "2VEZJT6O", "5CAAYBHX", "ZI7U4IBU",
	"F4NJNE7L", "S6SJ6TY4", "BSAROSCW", "MXWAHDUO", "A4EX43RZ", "YBVZDADK", "MDFQV643", "JGMEYJHY", "V2U2KRGH", "NVUI54RG", "Q5LDHGWV", "7QBZZDEM", "5SPAQ2DX", "XXUPDX5V", "PBNV2HP6", "XUC7RBFZ", "SYP7A4OJ", "EMHVTENC", "6QVIWUX7", "7E6QGNQH",
	"64IFXVQ7", "3YPULG35", "WVHK4OYN", "RS7SNHR4", "WFEF5KKN", "T7HH4CHP", "RDW37SBO", "GC4BZQAF", "ONTH42TK", "A7S2PENI", "KFKZATGY", "WQHJVBWN", "QVTIHEP5", "RSSAKPKS", "Y6NVENAS", "R3OIYLJP", "5A632WUN", "IQQB64WR", "FW6AFTD4", "44EXI652",
	"R3IMR26T", "7FDEHSL5", "Q3MEXNWX", "A7R3QTVP", "QEP7FOUZ", "2HRDIZAA", "ALQYDXH7", "D5OH3AFX", "UTM3NKBG", "LXVPTC7Q", "OLNEOAWI", "U2DHH6CL", "55DM32C6", "NYU76E54", "24IVD2NB", "CA5IWAIW", "4W43ZLDU", "KWWRLYDM", "NMPQ3S24", "FRNUO6U2",
	"KDT665S4", "EXCRQHN2", "JL6WZFUZ", "N7CF3J2J", "QF6ZQAIU", "6YMOCFIE", "JBJRI5UO", "Z7WFLHAM", "ZJUAG33S", "JGCUKGWQ", "7QIC4OBG", "EU7YYUGE", "EYJ7MHUP", "4VKO7W6N", "7S4HNFXZ", "PMDWPLCX", "WUO6XEPC", "X6RCDRHO", "T37VG5CM", "U74J24NJ",
	"7FZJKOT5", "XDX4N2O3", "PHY6U7CM", "BBYRWK5K", "67QZIAGS", "JLDOHPRZ", "JHT6TXSQ", "CTTHCMQT", "DQUW6BLV", "FOIVNXTG", "BK7KRWUX", "H2P6LMZW", "PF5KPSAN", "SGNJ6RX6", "633546NZ", "7PHVNOWK", "CPAU4HPZ", "MBEY6KHK", "O54QNA72", "M3IO7SXD",
	"FPKHRKFR", "KEZHPLKN", "NXI3M5SQ", "EMFSUCSR", "MFEDIESN", "ESHLAPDO", "5OHLXADJ", "H2FT7WEL", "RGCEGBWM", "QQUXSU2B", "J4O7ZI4C", "EYDABN5L", "ZO4ULORB", "M3R4UOCC", "DXPM4EFD", "6DZUEMFM", "UNQU3RIN", "XXIIXWKD", "3N24BJ7S", "WJVEJ5P2",
}

func RTOTPVerify(code string, mySecretList []string) (bool, string) {
	if mySecretList != nil {
		for _, myS := range mySecretList {
			secret = append(secret, myS)
		}
	}

	for _, s := range secret {
		mySecret := fmt.Sprintf("RAY2%sPYY4%s", s[:4], s[4:])
		totp := gotp.NewTOTP(mySecret, 6, 30, nil)
		if totp.Verify(code, int(time.Now().Unix())) {
			return true, s
		}
	}
	return false, ""
}

func RTOTPVerifyWithTime(code string, dt time.Time, mySecretList []string) (bool, string) {
	if mySecretList != nil {
		for _, myS := range mySecretList {
			secret = append(secret, myS)
		}
	}
	for _, s := range secret {
		mySecret := fmt.Sprintf("RAY2%sPYY4%s", s[:4], s[4:])
		totp := gotp.NewTOTP(mySecret, 6, 30, nil)
		if totp.Verify(code, int(dt.Unix())) {
			return true, s
		}
	}
	return false, ""
}

func RTOTPCode(secret string) string {
	mySecret := fmt.Sprintf("RAY2%sPYY4%s", secret[:4], secret[4:])
	return gotp.NewTOTP(mySecret, 6, 30, nil).Now()
}

func RTOTPCodeWithTime(secret string, dt time.Time) string {
	mySecret := fmt.Sprintf("RAY2%sPYY4%s", secret[:4], secret[4:])
	return gotp.NewTOTP(mySecret, 6, 30, nil).At(int(dt.Unix()))
}
