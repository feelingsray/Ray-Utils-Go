package rotp

import (
    "fmt"
    "time"
)

var secret  = []string{
    "rayray00","pyypyy00","raypyy86","raypyy89","ray89pyy",
    "2d7bk8lg", "l4hdi1eu", "dmtdthbz", "ioxbmvx6", "kjqwyt2l", "10ml50ho", "1ws0t0ml", "akubfgkl", "08h7dbpd", "7qbnry2g",
    "k96pa6pa", "yjfgy21k", "v8n0951o", "kgei6hp7", "l0ugvhbo", "4gf974gi", "mm5punlq", "woyg7pj5", "osy3jerr", "v19iyo7p",
    "r3vjv6tb", "ra8wti17", "290kdesv", "krjgncxs", "ieiufg7c", "uxs8eoed", "ej477uhg", "jsuj1ccz", "izwdqwk6", "ecl4q84q",
    "gl1ch1k6", "jkxj5ris", "dc8juuw1", "ur11a40u", "szfhbqjt", "mec662vn", "roeic7rj", "pbdzy6gy", "myy9sdt5", "dw953i60",
    "ifebhmm3", "mru8o7zu", "gz4edilk", "at4wl4fc", "w4xw5vkl", "6zs0f31b", "rnoenr2m", "1xz6wvbb", "jdg208qp", "elbdufn7",
    "w69tc6b5", "84ba3fgr", "0h8v1hrz", "fj2btpzy", "m0ptarz8", "gki3xdyk", "i7l6crfx", "6v9bkgt8", "m5oszmbb", "0vohodll",
    "kkb2m4td", "k0c04sg3", "pl7ei7ry", "gwpv2pbd", "50b9t5wa", "jwkwmv8r", "zdkhsg3z", "q8kfwssp", "yes5e574", "414rlnlr",
    "6vayetl8", "k365pk0n", "czuykkm4", "2abiv5j7", "qjtq89rj", "5xnap694", "3uawp5c2", "xvauoncd", "8czfjevl", "wokgh7wv",
    "crqnfdfe", "377qozh8", "4xfpz7ne", "1mfp2thc", "85dogw4k", "yomdewm7", "282cs387", "dd41wpg8", "p8zzii0i", "fxiuv8bg",
    "k9alady7", "6am5kl8y", "o8h5j5ph", "57av6bja", "457dwhzd", "6ueo3vna", "fuaw34t3", "e3ap7hnd", "ayfnleec", "kdqlza6o",
    "thop9i7s", "a9g9c6af", "xbo3ezy3", "fedbz9qs", "hdxjdfwq", "h0lmnk99", "a59y1lv3", "8cp2jkt8", "yi108ok0", "n18yi9yb",
    "lc8ynj9t", "qdh30y9l", "q7p1fvck", "uuaahxcb", "l3jrsriy", "e497j2g8", "zzknnubf", "k1ngcz2c", "fzjqaiig", "xzkdc6zs",
    "e555a3jb", "0jmcvryy", "2zl5x3ph", "sox1bgav", "py50cal5", "lpg2vg82", "xqxgz2ik", "xbl9uxt0", "8kah8cte", "jnkszfo0",
    "eqjbkxi2", "5d7ipzeq", "95tlmm5r", "sohvwudy", "d7xi5x5h", "2o6l9s1h", "oeofs8ae", "3rxt7fmj", "32kwflej", "a82aya72",
    "2bfyfcd3", "poedu6vr", "rq37myu2", "0xiwq0xd", "h2kl2iia", "565lm15s", "87nunm1e", "q726ku5t", "3swltfiz", "9uu8moc0",
    "nmdaw80s", "khrncj04", "7zikpt2f", "6kq657bd", "d69myshp", "90r34ipk", "8h8bacdp", "9lkk7ux0", "3t3ja6pn", "aw5xascr",
    "7k84l2mq", "ywh11bcs", "h8pq39uz", "qxgt2tep", "cebcjlnw", "xecxm2qe", "5aqrva9m", "vjskdy20", "7xlgnjc7", "oocueqy6",
    "mshao7ik", "xwbd9ro1", "nnfgtyle", "vpqmfmsg", "tbzyl8mg", "u7s4a3xw", "opqmdip8", "32286liy", "z7o13qqw", "owuiixzd",
    "6b0x1lbg", "d8oriplz", "jl2qfd3i", "abr92ika", "pf7rkl29", "0yj5l6cc", "6mc5f6z7", "sefkm03h", "orkij4wa", "sgxpl276",
    "4knihdt2", "qotbpnr8", "ws8tksct", "ow14bs8z", "3dr2og9u", "q6lxc6yc", "vbd6iwig", "b7xp5rt6", "7qljdwnc", "h5w2pisn",
    "pcncng9d", "kfh0vpkc", "mionqtqu", "qfup0iuo", "wle8vytb", "vpxfcdnc", "ly81hpdo", "5kwg4yjj", "wt9wo1k9", "ws92qjbp",
    "p4gbzjd8", "51rbxf82", "ag1vlue2", "p5z2c4d1", "kbmi1ors", "1j9tgjpk", "72iujawt", "hchtthyi", "g87ogcna", "7ftyxy01",
    "66k2g71v", "i74zsxh3", "cs9is933", "sq460md2", "un590jdb", "vyb8qo2o", "2occrej6", "8ssmvw5k", "vtqb8mbr", "drldqqrm",
    "tcnssuvs", "f2y9rq1k", "q5zptc8p", "t825kc2e", "wxypa2gs", "1591niji", "83c49piv", "i3n1ffwk", "ufzhrcya", "e8o44hfb",
    "rhonnrlq", "4cu9n96q", "i219nd8t", "6ju4xswt", "2plg5yi4", "i3zuomgb", "3l0rgvsh", "quheu737", "g8cpz18l", "r92to87z",
    "1qh7wm8z", "wjkit60f", "567uv44o", "7d64kyae", "mafjm6ov", "2qf43kl2", "pquqigr2", "vgl4d2hq", "3fm7ewny", "duhdmy28",
    "kjvorryp", "y6ma3g19", "r27j2l8h", "w6qldm8i", "zdgoxa4m", "2xkcidyg", "bpwow6lf", "7bl1swbg", "rh2a447d", "fccg1arx",
    "aw0kt9qm", "g6tzk4cb", "gsaxge4v", "0v6f8i2k", "3smwbn98", "rwqw1z1y", "utcvsgf4", "mw1cd4no", "87q210yr", "yi8x2jyi",
    "3x80fc0i", "5u3967za", "nkue6gge", "2bj2l18z", "tp4ecqhn", "oo9kzegd", "olr8dupf", "4odj528q", "j3j6tqzi", "wee2pbae",
    "mx5bm5m8", "9545hfei", "1m69ncvo", "ap5f5p33", "j6vzndab", "3nd16lmv", "haqk4ri1", "j9u354du", "iz8n1q6h", "8o5h4hj7",
    "03wemwup", "0f9i6onb", "sx9ggm22", "skr07ipb", "ogyvswcy", "nwnpvkbn", "486qbc18", "adrnvs70", "omffn9n2", "nvqe8dwi",
    "ihvhdyzl", "sgx1o7rj", "io9268s4", "538m2sws", "hpv0pfv0", "ahmb5jl2", "qgy8y8ht", "p0ep4g41", "18jj2mwz", "f7ir4w6l",
    "5aiz5m80", "talwn5es", "1efyon3l", "1qbg21l8", "wcmpjen0", "8jem3ger", "jlzu7ozs", "o2frftme", "92jfjpqr", "r1ni302f",
    "9j9ry285", "pclrzzj7", "feexxk9g", "ms954thg", "1cocqoch", "2kw04y6c", "cwlvkon8", "1tcmlzwi", "5477gore", "sabofpq8",
    "axz3sruu", "hsq7siqf", "mi90etxs", "ksd61a3g", "9csa3yl1", "563g6m47", "c5zbmj1k", "69wm2hf8", "a67ol50j", "yyshuopm",
    "g453awwi", "9e0zdfd2", "stt8psaw", "rs2lkefa", "vpwlhjwd", "hplpcwmx", "wrcj4bma", "lwm5m3w0", "yty4fdt8", "m6h88mpu",
    "f1xucm6x", "ofj2t7hf", "bjxyc1sz", "v11nyag7", "5q1fg8gg", "j20i01b1", "oibj7g9f", "hyy11c79", "juw8ittt", "o4xt3cfz",
    "t0o2v8pg", "h0067fvs", "v94in9ry", "1s0nbn9l", "94e6oud5", "cc13eyho", "smao67ob", "m3zqh596", "nxxmymah", "uj1onvrd",
    "qc6h0y5x", "4pa23civ", "j1b0d7hg", "5o8zs58t", "2ju4hur4", "cmtr5o6f", "rs7w96r2", "mh4l99si", "zfq3rsrg", "gu2tj3ch",
    "mcvn9dh3", "8hzwh9d9", "wraq5xy8", "9o7o0egz", "8ovo89a2", "l63wkuqo", "qb4upjex", "8oa585ta", "0a1tcfk7", "ilbazw2b",
    "ep5r8mpq", "e2yfh45t", "ad06dmww", "abe76jaj", "1274kjgm", "qqtgpq8h", "pvabl8vm", "vj1kley1", "fnfxeey9", "vaks554l",
    "37gf7g7k", "8zcw3gor", "69asj7bk", "4atuzham", "prmxz3iv", "l7xf6lqh", "c3q1qtxv", "rzbfeqfp", "aiw55etb", "brj7aqbl",
    "ec8j3cqy", "xsb20ujo", "d7az0lqt", "uajq37a0", "0difkgxn", "jx6l10nx", "s0q54qr1", "sw36g1yl", "nl0xwehw", "szpj0jro",
    "b5ta3dsl", "m1h0qtxm", "ilhhvwk4", "vzc03dyy", "k6ln1hl8", "ke7zmemz", "urqo1dku", "e1skgptf", "hbw72j8i", "h8h8z1gv",
    "n9iwgl89", "ffdj1lkm", "z5rq43d9", "72qo2apr", "or9bc1hy", "onig6jch", "hrprd9cj", "9rfhi9t1", "8vr94tkz", "5hw3fnx2",
    "6jz6ghc6", "k7y4dyw7", "32sq8715", "6f4wyu10", "fmchv0vx", "ny2mfw2z", "ewfmmihi", "1w7h7ypo", "0xoc6gcv", "doukgliu",
    "g7cpqy1a", "4xr1eo3c", "jt16mmsr", "kwrbj015", "owd939lf", "t45813q0", "wsawcznk", "r8ic7p7o", "pm1swdrt", "fnjxqt1w",
    "a19xtsuy", "2ie1xsgp", "r8t3qjwt", "gz2awgrv", "hxq3gj4p", "55cnoozf", "1als9r02", "04hron2e", "uv2royc2", "x247bkm8",
    "j4oflt3g", "b7ofi2e1", "jq60w46q", "68le4fwx", "b89mqyi5", "qd97rxq0", "d50eae4e", "go8n0nn9", "xw69c4hm", "8bcdccgg",
    "3qas3l1d", "y3nvcqhh", "hkj1b67c", "8rr8isj6", "nrhbyh09", "u1uxct32", "z1pikvez", "rqszv9dy", "ak01en1p", "zatwic1l",
}


func RTOTPVerify(code string) (bool,string) {
    //fmt.Println(len(secret))
    for _,s := range secret {
        mySecret := fmt.Sprintf("ray%s1989%s02",s[:4],s[4:])
        //fmt.Println(mySecret)
        if NewTOTP([]byte(mySecret),6,30).Verify(code) {
            return true,s
        }
    }
    return false,""
}

func RTOTPVerifyWithTime(code string,dt time.Time) (bool,string) {
    for _,s := range secret {
        mySecret := fmt.Sprintf("ray%s1989%s02",s[:4],s[4:])
        if NewTOTP([]byte(mySecret),6,30).VerifyWithTime(code,dt) {
            return true,s
        }
    }
    return false,""
}

func RTOTPCode(secret string) string {
    mySecret := fmt.Sprintf("ray%s1989%s02",secret[:4],secret[4:])
    return NewTOTP([]byte(mySecret),6,30).Now()
}

func RTOTPCodeWithTime(secret string,dt time.Time) string {
    mySecret := fmt.Sprintf("ray%s1989%s02",secret[:4],secret[4:])
    return NewTOTP([]byte(mySecret),6,30).Time(dt)
}