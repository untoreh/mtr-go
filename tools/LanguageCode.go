package tools

import (
	"log"
)

var Langs = map[string]map[string]string{
	"auto": map[string]string{
		"name":       "auto",
		"native":     "auto",
		"iso639-1":   "auto",
		"iso639-2/t": "auto",
		"iso639-2/b": "auto",
		"iso639-3":   "auto",
		// "bing": "auto", // only from api
		"yandex": "auto",
		"google": "auto",
		// "multillect": "auto", // only from api
		"frengly": "auto",
		"sdl":     "auto",
		"systran": "auto",
		"promt":   "auto",
	},
	"ab": map[string]string{
		"name":       "abkhaz",
		"native":     "аҧсуа бызшәа, аҧсшәа",
		"iso639-1":   "ab",
		"iso639-2/t": "abk",
		"iso639-2/b": "abk",
		"iso639-3":   "abk",
	},
	"aa": map[string]string{
		"name":       "afar",
		"native":     "afaraf",
		"iso639-1":   "aa",
		"iso639-2/t": "aar",
		"iso639-2/b": "aar",
		"iso639-3":   "aar",
	},
	"af": map[string]string{
		"name":       "afrikaans",
		"native":     "afrikaans",
		"iso639-1":   "af",
		"iso639-2/t": "afr",
		"iso639-2/b": "afr",
		"iso639-3":   "afr",
	},
	"ak": map[string]string{
		"name":       "akan",
		"native":     "akan",
		"iso639-1":   "ak",
		"iso639-2/t": "aka",
		"iso639-2/b": "aka",
		"iso639-3":   "aka + 2",
	},
	"sq": map[string]string{
		"name":       "albanian",
		"native":     "shqip",
		"iso639-1":   "sq",
		"iso639-2/t": "sqi",
		"iso639-2/b": "alb",
		"iso639-3":   "sqi + 4",
	},
	"am": map[string]string{
		"name":       "amharic",
		"native":     "አማርኛ",
		"iso639-1":   "am",
		"iso639-2/t": "amh",
		"iso639-2/b": "amh",
		"iso639-3":   "amh",
	},
	"ar": map[string]string{
		"name":       "arabic",
		"native":     "العربية",
		"iso639-1":   "ar",
		"iso639-2/t": "ara",
		"iso639-2/b": "ara",
		"iso639-3":   "ara + 30",
	},
	"an": map[string]string{
		"name":       "aragonese",
		"native":     "aragonés",
		"iso639-1":   "an",
		"iso639-2/t": "arg",
		"iso639-2/b": "arg",
		"iso639-3":   "arg",
	},
	"hy": map[string]string{
		"name":       "armenian",
		"native":     "հայերեն",
		"iso639-1":   "hy",
		"iso639-2/t": "hye",
		"iso639-2/b": "arm",
		"iso639-3":   "hye",
	},
	"as": map[string]string{
		"name":       "assamese",
		"native":     "অসমীয়া",
		"iso639-1":   "as",
		"iso639-2/t": "asm",
		"iso639-2/b": "asm",
		"iso639-3":   "asm",
	},
	"av": map[string]string{
		"name":       "avaric",
		"native":     "авар мацӏ, магӏарул мацӏ",
		"iso639-1":   "av",
		"iso639-2/t": "ava",
		"iso639-2/b": "ava",
		"iso639-3":   "ava",
	},
	"ae": map[string]string{
		"name":       "avestan",
		"native":     "avesta",
		"iso639-1":   "ae",
		"iso639-2/t": "ave",
		"iso639-2/b": "ave",
		"iso639-3":   "ave",
	},
	"ay": map[string]string{
		"name":       "aymara",
		"native":     "aymar aru",
		"iso639-1":   "ay",
		"iso639-2/t": "aym",
		"iso639-2/b": "aym",
		"iso639-3":   "aym + 2",
	},
	"az": map[string]string{
		"name":       "azerbaijani",
		"native":     "azərbaycan dili",
		"iso639-1":   "az",
		"iso639-2/t": "aze",
		"iso639-2/b": "aze",
		"iso639-3":   "aze + 2",
	},
	"bm": map[string]string{
		"name":       "bambara",
		"native":     "bamanankan",
		"iso639-1":   "bm",
		"iso639-2/t": "bam",
		"iso639-2/b": "bam",
		"iso639-3":   "bam",
	},
	"ba": map[string]string{
		"name":       "bashkir",
		"native":     "башҡорт теле",
		"iso639-1":   "ba",
		"iso639-2/t": "bak",
		"iso639-2/b": "bak",
		"iso639-3":   "bak",
	},
	"eu": map[string]string{
		"name":       "basque",
		"native":     "euskara, euskera",
		"iso639-1":   "eu",
		"iso639-2/t": "eus",
		"iso639-2/b": "baq",
		"iso639-3":   "eus",
	},
	"be": map[string]string{
		"name":       "belarusian",
		"native":     "беларуская мова",
		"iso639-1":   "be",
		"iso639-2/t": "bel",
		"iso639-2/b": "bel",
		"iso639-3":   "bel",
	},
	"bn": map[string]string{
		"name":       "bengali, bangla",
		"native":     "বাংলা",
		"iso639-1":   "bn",
		"iso639-2/t": "ben",
		"iso639-2/b": "ben",
		"iso639-3":   "ben",
		"bing":       "bn-BD",
	},
	"bh": map[string]string{
		"name":       "bihari",
		"native":     "भोजपुरी",
		"iso639-1":   "bh",
		"iso639-2/t": "bih",
		"iso639-2/b": "bih",
		"iso639-3":   "bih",
	},
	"bi": map[string]string{
		"name":       "bislama",
		"native":     "bislama",
		"iso639-1":   "bi",
		"iso639-2/t": "bis",
		"iso639-2/b": "bis",
		"iso639-3":   "bis",
	},
	"bs": map[string]string{
		"name":       "bosnian",
		"native":     "bosanski jezik",
		"iso639-1":   "bs",
		"iso639-2/t": "bos",
		"iso639-2/b": "bos",
		"iso639-3":   "bos",
	},
	"br": map[string]string{
		"name":       "breton",
		"native":     "brezhoneg",
		"iso639-1":   "br",
		"iso639-2/t": "bre",
		"iso639-2/b": "bre",
		"iso639-3":   "bre",
	},
	"bg": map[string]string{
		"name":       "bulgarian",
		"native":     "български език",
		"iso639-1":   "bg",
		"iso639-2/t": "bul",
		"iso639-2/b": "bul",
		"iso639-3":   "bul",
	},
	"my": map[string]string{
		"name":       "burmese",
		"native":     "ဗမာစာ",
		"iso639-1":   "my",
		"iso639-2/t": "mya",
		"iso639-2/b": "bur",
		"iso639-3":   "mya",
	},
	"ca": map[string]string{
		"name":       "catalan, valencian",
		"native":     "català, valencià",
		"iso639-1":   "ca",
		"iso639-2/t": "cat",
		"iso639-2/b": "cat",
		"iso639-3":   "cat",
	},
	"ceb": map[string]string{
		"name":       "cebuano",
		"native":     "cebuano",
		"iso639-1":   "ceb",
		"iso639-2/t": "ceb",
		"iso639-2/b": "ceb",
		"iso639-3":   "ceb",
	},
	"ch": map[string]string{
		"name":       "chamorro",
		"native":     "chamoru",
		"iso639-1":   "ch",
		"iso639-2/t": "cha",
		"iso639-2/b": "cha",
		"iso639-3":   "cha",
	},
	"ce": map[string]string{
		"name":       "chechen",
		"native":     "нохчийн мотт",
		"iso639-1":   "ce",
		"iso639-2/t": "che",
		"iso639-2/b": "che",
		"iso639-3":   "che",
	},
	"ny": map[string]string{
		"name":       "chichewa, chewa, nyanja",
		"native":     "chicheŵa, chinyanja",
		"iso639-1":   "ny",
		"iso639-2/t": "nya",
		"iso639-2/b": "nya",
		"iso639-3":   "nya",
	},
	//"zh": map[string]string{
	//	"name":       "chinese",
	//	"native":     "中文 (zhōngwén), 汉语, 漢語",
	//	"iso639-1":   "zh",
	//	"iso639-2/t": "zho",
	//	"iso639-2/b": "chi",
	//	"iso639-3":   "zho + 13",
	//	"google":     "zh",
	//	"treu": "zh-CN",
	//},

	"cv": map[string]string{
		"name":       "chuvash",
		"native":     "чӑваш чӗлхи",
		"iso639-1":   "cv",
		"iso639-2/t": "chv",
		"iso639-2/b": "chv",
		"iso639-3":   "chv",
	},
	"kw": map[string]string{
		"name":       "cornish",
		"native":     "kernewek",
		"iso639-1":   "kw",
		"iso639-2/t": "cor",
		"iso639-2/b": "cor",
		"iso639-3":   "cor",
	},
	"co": map[string]string{
		"name":       "corsican",
		"native":     "corsu, lingua corsa",
		"iso639-1":   "co",
		"iso639-2/t": "cos",
		"iso639-2/b": "cos",
		"iso639-3":   "cos",
	},
	"cr": map[string]string{
		"name":       "cree",
		"native":     "ᓀᐦᐃᔭᐍᐏᐣ",
		"iso639-1":   "cr",
		"iso639-2/t": "cre",
		"iso639-2/b": "cre",
		"iso639-3":   "cre + 6",
	},
	"hr": map[string]string{
		"name":       "croatian",
		"native":     "hrvatski jezik",
		"iso639-1":   "hr",
		"iso639-2/t": "hrv",
		"iso639-2/b": "hrv",
		"iso639-3":   "hrv",
	},
	"cs": map[string]string{
		"name":       "czech",
		"native":     "čeština, český jazyk",
		"iso639-1":   "cs",
		"iso639-2/t": "ces",
		"iso639-2/b": "cze",
		"iso639-3":   "ces",
	},
	"da": map[string]string{
		"name":       "danish",
		"native":     "dansk",
		"iso639-1":   "da",
		"iso639-2/t": "dan",
		"iso639-2/b": "dan",
		"iso639-3":   "dan",
	},
	"dv": map[string]string{
		"name":       "divehi, dhivehi, maldivian",
		"native":     "ދިވެހި",
		"iso639-1":   "dv",
		"iso639-2/t": "div",
		"iso639-2/b": "div",
		"iso639-3":   "div",
	},
	"nl": map[string]string{
		"name":       "dutch",
		"native":     "nederlands, vlaams",
		"iso639-1":   "nl",
		"iso639-2/t": "nld",
		"iso639-2/b": "dut",
		"iso639-3":   "nld",
	},
	"dz": map[string]string{
		"name":       "dzongkha",
		"native":     "རྫོང་ཁ",
		"iso639-1":   "dz",
		"iso639-2/t": "dzo",
		"iso639-2/b": "dzo",
		"iso639-3":   "dzo",
	},
	"en": map[string]string{
		"name":       "english",
		"native":     "english",
		"iso639-1":   "en",
		"iso639-2/t": "eng",
		"iso639-2/b": "eng",
		"iso639-3":   "eng",
	},
	"eo": map[string]string{
		"name":       "esperanto",
		"native":     "esperanto",
		"iso639-1":   "eo",
		"iso639-2/t": "epo",
		"iso639-2/b": "epo",
		"iso639-3":   "epo",
	},
	"et": map[string]string{
		"name":       "estonian",
		"native":     "eesti, eesti keel",
		"iso639-1":   "et",
		"iso639-2/t": "est",
		"iso639-2/b": "est",
		"iso639-3":   "est + 2",
	},
	"ee": map[string]string{
		"name":       "ewe",
		"native":     "eʋegbe",
		"iso639-1":   "ee",
		"iso639-2/t": "ewe",
		"iso639-2/b": "ewe",
		"iso639-3":   "ewe",
	},
	"fo": map[string]string{
		"name":       "faroese",
		"native":     "føroyskt",
		"iso639-1":   "fo",
		"iso639-2/t": "fao",
		"iso639-2/b": "fao",
		"iso639-3":   "fao",
	},
	"fj": map[string]string{
		"name":       "fijian",
		"native":     "vosa vakaviti",
		"iso639-1":   "fj",
		"iso639-2/t": "fij",
		"iso639-2/b": "fij",
		"iso639-3":   "fij",
	},
	"fi": map[string]string{
		"name":       "finnish",
		"native":     "suomi, suomen kieli",
		"iso639-1":   "fi",
		"iso639-2/t": "fin",
		"iso639-2/b": "fin",
		"iso639-3":   "fin",
	},
	"fr": map[string]string{
		"name":       "french",
		"native":     "français, langue française",
		"iso639-1":   "fr",
		"iso639-2/t": "fra",
		"iso639-2/b": "fre",
		"iso639-3":   "fra",
	},
	"ff": map[string]string{
		"name":       "fula, fulah, pulaar, pular",
		"native":     "fulfulde, pulaar, pular",
		"iso639-1":   "ff",
		"iso639-2/t": "ful",
		"iso639-2/b": "ful",
		"iso639-3":   "ful + 9",
	},
	"gl": map[string]string{
		"name":       "galician",
		"native":     "galego",
		"iso639-1":   "gl",
		"iso639-2/t": "glg",
		"iso639-2/b": "glg",
		"iso639-3":   "glg",
	},
	"ka": map[string]string{
		"name":       "georgian",
		"native":     "ქართული",
		"iso639-1":   "ka",
		"iso639-2/t": "kat",
		"iso639-2/b": "geo",
		"iso639-3":   "kat",
	},
	"de": map[string]string{
		"name":       "german",
		"native":     "deutsch",
		"iso639-1":   "de",
		"iso639-2/t": "deu",
		"iso639-2/b": "ger",
		"iso639-3":   "deu",
	},
	"el": map[string]string{
		"name":       "greek (modern)",
		"native":     "ελληνικά",
		"iso639-1":   "el",
		"iso639-2/t": "ell",
		"iso639-2/b": "gre",
		"iso639-3":   "ell",
	},
	"gn": map[string]string{
		"name":       "guaraní",
		"native":     "avañe\"ẽ",
		"iso639-1":   "gn",
		"iso639-2/t": "grn",
		"iso639-2/b": "grn",
		"iso639-3":   "grn + 5",
	},
	"gu": map[string]string{
		"name":       "gujarati",
		"native":     "ગુજરાતી",
		"iso639-1":   "gu",
		"iso639-2/t": "guj",
		"iso639-2/b": "guj",
		"iso639-3":   "guj",
	},
	"ht": map[string]string{
		"name":       "haitian, haitian creole",
		"native":     "kreyòl ayisyen",
		"iso639-1":   "ht",
		"iso639-2/t": "hat",
		"iso639-2/b": "hat",
		"iso639-3":   "hat",
	},
	"ha": map[string]string{
		"name":       "hausa",
		"native":     "(hausa) هَوُسَ",
		"iso639-1":   "ha",
		"iso639-2/t": "hau",
		"iso639-2/b": "hau",
		"iso639-3":   "hau",
	},
	"he": map[string]string{
		"name":       "hebrew (modern)",
		"native":     "עברית",
		"iso639-1":   "he",
		"iso639-2/t": "heb",
		"iso639-2/b": "heb",
		"iso639-3":   "heb",
	},
	"hz": map[string]string{
		"name":       "herero",
		"native":     "otjiherero",
		"iso639-1":   "hz",
		"iso639-2/t": "her",
		"iso639-2/b": "her",
		"iso639-3":   "her",
	},
	"hi": map[string]string{
		"name":       "hindi",
		"native":     "हिन्दी, हिंदी",
		"iso639-1":   "hi",
		"iso639-2/t": "hin",
		"iso639-2/b": "hin",
		"iso639-3":   "hin",
	},
	"ho": map[string]string{
		"name":       "hiri motu",
		"native":     "hiri motu",
		"iso639-1":   "ho",
		"iso639-2/t": "hmo",
		"iso639-2/b": "hmo",
		"iso639-3":   "hmo",
	},
	"hu": map[string]string{
		"name":       "hungarian",
		"native":     "magyar",
		"iso639-1":   "hu",
		"iso639-2/t": "hun",
		"iso639-2/b": "hun",
		"iso639-3":   "hun",
	},
	"ia": map[string]string{
		"name":       "interlingua",
		"native":     "interlingua",
		"iso639-1":   "ia",
		"iso639-2/t": "ina",
		"iso639-2/b": "ina",
		"iso639-3":   "ina",
	},
	"id": map[string]string{
		"name":       "indonesian",
		"native":     "bahasa indonesia",
		"iso639-1":   "id",
		"iso639-2/t": "ind",
		"iso639-2/b": "ind",
		"iso639-3":   "ind",
	},
	"ie": map[string]string{
		"name":       "interlingue",
		"native":     "originally called occidental; then interlingue after WWII",
		"iso639-1":   "ie",
		"iso639-2/t": "ile",
		"iso639-2/b": "ile",
		"iso639-3":   "ile",
	},
	"ga": map[string]string{
		"name":       "irish",
		"native":     "gaeilge",
		"iso639-1":   "ga",
		"iso639-2/t": "gle",
		"iso639-2/b": "gle",
		"iso639-3":   "gle",
	},
	"ig": map[string]string{
		"name":       "igbo",
		"native":     "asụsụ igbo",
		"iso639-1":   "ig",
		"iso639-2/t": "ibo",
		"iso639-2/b": "ibo",
		"iso639-3":   "ibo",
	},
	"ik": map[string]string{
		"name":       "inupiaq",
		"native":     "iñupiaq, iñupiatun",
		"iso639-1":   "ik",
		"iso639-2/t": "ipk",
		"iso639-2/b": "ipk",
		"iso639-3":   "ipk + 2",
	},
	"io": map[string]string{
		"name":       "ido",
		"native":     "ido",
		"iso639-1":   "io",
		"iso639-2/t": "ido",
		"iso639-2/b": "ido",
		"iso639-3":   "ido",
	},
	"is": map[string]string{
		"name":       "icelandic",
		"native":     "íslenska",
		"iso639-1":   "is",
		"iso639-2/t": "isl",
		"iso639-2/b": "ice",
		"iso639-3":   "isl",
	},
	"it": map[string]string{
		"name":       "italian",
		"native":     "italiano",
		"iso639-1":   "it",
		"iso639-2/t": "ita",
		"iso639-2/b": "ita",
		"iso639-3":   "ita",
	},
	"iu": map[string]string{
		"name":       "inuktitut",
		"native":     "ᐃᓄᒃᑎᑐᑦ",
		"iso639-1":   "iu",
		"iso639-2/t": "iku",
		"iso639-2/b": "iku",
		"iso639-3":   "iku + 2",
	},
	"ja": map[string]string{
		"name":       "japanese",
		"native":     "日本語 (にほんご)",
		"iso639-1":   "ja",
		"iso639-2/t": "jpn",
		"iso639-2/b": "jpn",
		"iso639-3":   "jpn",
	},
	"jv": map[string]string{
		"name":       "javanese",
		"native":     "basa jawa",
		"iso639-1":   "jv",
		"iso639-2/t": "jav",
		"iso639-2/b": "jav",
		"iso639-3":   "jav",
	},
	"kl": map[string]string{
		"name":       "kalaallisut, greenlandic",
		"native":     "kalaallisut, kalaallit oqaasii",
		"iso639-1":   "kl",
		"iso639-2/t": "kal",
		"iso639-2/b": "kal",
		"iso639-3":   "kal",
	},
	"kn": map[string]string{
		"name":       "kannada",
		"native":     "ಕನ್ನಡ",
		"iso639-1":   "kn",
		"iso639-2/t": "kan",
		"iso639-2/b": "kan",
		"iso639-3":   "kan",
	},
	"kr": map[string]string{
		"name":       "kanuri",
		"native":     "kanuri",
		"iso639-1":   "kr",
		"iso639-2/t": "kau",
		"iso639-2/b": "kau",
		"iso639-3":   "kau + 3",
	},
	"ks": map[string]string{
		"name":       "kashmiri",
		"native":     "कश्मीरी, كشميري‎",
		"iso639-1":   "ks",
		"iso639-2/t": "kas",
		"iso639-2/b": "kas",
		"iso639-3":   "kas",
	},
	"kk": map[string]string{
		"name":       "kazakh",
		"native":     "қазақ тілі",
		"iso639-1":   "kk",
		"iso639-2/t": "kaz",
		"iso639-2/b": "kaz",
		"iso639-3":   "kaz",
	},
	"km": map[string]string{
		"name":       "khmer",
		"native":     "ខ្មែរ, ខេមរភាសា, ភាសាខ្មែរ",
		"iso639-1":   "km",
		"iso639-2/t": "khm",
		"iso639-2/b": "khm",
		"iso639-3":   "khm",
	},
	"ki": map[string]string{
		"name":       "kikuyu, gikuyu",
		"native":     "gĩkũyũ",
		"iso639-1":   "ki",
		"iso639-2/t": "kik",
		"iso639-2/b": "kik",
		"iso639-3":   "kik",
	},
	"rw": map[string]string{
		"name":       "kinyarwanda",
		"native":     "ikinyarwanda",
		"iso639-1":   "rw",
		"iso639-2/t": "kin",
		"iso639-2/b": "kin",
		"iso639-3":   "kin",
	},
	"ky": map[string]string{
		"name":       "kyrgyz",
		"native":     "кыргызча, кыргыз тили",
		"iso639-1":   "ky",
		"iso639-2/t": "kir",
		"iso639-2/b": "kir",
		"iso639-3":   "kir",
	},
	"kv": map[string]string{
		"name":       "komi",
		"native":     "коми кыв",
		"iso639-1":   "kv",
		"iso639-2/t": "kom",
		"iso639-2/b": "kom",
		"iso639-3":   "kom + 2",
	},
	"kg": map[string]string{
		"name":       "kongo",
		"native":     "kikongo",
		"iso639-1":   "kg",
		"iso639-2/t": "kon",
		"iso639-2/b": "kon",
		"iso639-3":   "kon + 3",
	},
	"ko": map[string]string{
		"name":       "korean",
		"native":     "한국어, 조선어",
		"iso639-1":   "ko",
		"iso639-2/t": "kor",
		"iso639-2/b": "kor",
		"iso639-3":   "kor",
	},
	"ku": map[string]string{
		"name":       "kurdish",
		"native":     "kurdî, كوردی‎",
		"iso639-1":   "ku",
		"iso639-2/t": "kur",
		"iso639-2/b": "kur",
		"iso639-3":   "kur + 3",
	},
	"kj": map[string]string{
		"name":       "kwanyama, kuanyama",
		"native":     "kuanyama",
		"iso639-1":   "kj",
		"iso639-2/t": "kua",
		"iso639-2/b": "kua",
		"iso639-3":   "kua",
	},
	"la": map[string]string{
		"name":       "latin",
		"native":     "latine, lingua latina",
		"iso639-1":   "la",
		"iso639-2/t": "lat",
		"iso639-2/b": "lat",
		"iso639-3":   "lat",
	},
	"lb": map[string]string{
		"name":       "luxembourgish, letzeburgesch",
		"native":     "lëtzebuergesch",
		"iso639-1":   "lb",
		"iso639-2/t": "ltz",
		"iso639-2/b": "ltz",
		"iso639-3":   "ltz",
	},
	"lg": map[string]string{
		"name":       "ganda",
		"native":     "luganda",
		"iso639-1":   "lg",
		"iso639-2/t": "lug",
		"iso639-2/b": "lug",
		"iso639-3":   "lug",
	},
	"li": map[string]string{
		"name":       "limburgish, limburgan, limburger",
		"native":     "limburgs",
		"iso639-1":   "li",
		"iso639-2/t": "lim",
		"iso639-2/b": "lim",
		"iso639-3":   "lim",
	},
	"ln": map[string]string{
		"name":       "lingala",
		"native":     "lingála",
		"iso639-1":   "ln",
		"iso639-2/t": "lin",
		"iso639-2/b": "lin",
		"iso639-3":   "lin",
	},
	"lo": map[string]string{
		"name":       "lao",
		"native":     "ພາສາລາວ",
		"iso639-1":   "lo",
		"iso639-2/t": "lao",
		"iso639-2/b": "lao",
		"iso639-3":   "lao",
	},
	"lt": map[string]string{
		"name":       "lithuanian",
		"native":     "lietuvių kalba",
		"iso639-1":   "lt",
		"iso639-2/t": "lit",
		"iso639-2/b": "lit",
		"iso639-3":   "lit",
	},
	"lu": map[string]string{
		"name":       "luba-katanga",
		"native":     "tshiluba",
		"iso639-1":   "lu",
		"iso639-2/t": "lub",
		"iso639-2/b": "lub",
		"iso639-3":   "lub",
	},
	"lv": map[string]string{
		"name":       "latvian",
		"native":     "latviešu valoda",
		"iso639-1":   "lv",
		"iso639-2/t": "lav",
		"iso639-2/b": "lav",
		"iso639-3":   "lav + 2",
	},
	"gv": map[string]string{
		"name":       "manx",
		"native":     "gaelg, gailck",
		"iso639-1":   "gv",
		"iso639-2/t": "glv",
		"iso639-2/b": "glv",
		"iso639-3":   "glv",
	},
	"mk": map[string]string{
		"name":       "macedonian",
		"native":     "македонски јазик",
		"iso639-1":   "mk",
		"iso639-2/t": "mkd",
		"iso639-2/b": "mac",
		"iso639-3":   "mkd",
	},
	"mg": map[string]string{
		"name":       "malagasy",
		"native":     "fiteny malagasy",
		"iso639-1":   "mg",
		"iso639-2/t": "mlg",
		"iso639-2/b": "mlg",
		"iso639-3":   "mlg + 10",
	},
	"ms": map[string]string{
		"name":       "malay",
		"native":     "bahasa melayu, بهاس ملايو‎",
		"iso639-1":   "ms",
		"iso639-2/t": "msa",
		"iso639-2/b": "may",
		"iso639-3":   "msa + 13",
	},
	"ml": map[string]string{
		"name":       "malayalam",
		"native":     "മലയാളം",
		"iso639-1":   "ml",
		"iso639-2/t": "mal",
		"iso639-2/b": "mal",
		"iso639-3":   "mal",
	},
	"mt": map[string]string{
		"name":       "maltese",
		"native":     "malti",
		"iso639-1":   "mt",
		"iso639-2/t": "mlt",
		"iso639-2/b": "mlt",
		"iso639-3":   "mlt",
	},
	"mi": map[string]string{
		"name":       "māori",
		"native":     "te reo māori",
		"iso639-1":   "mi",
		"iso639-2/t": "mri",
		"iso639-2/b": "mao",
		"iso639-3":   "mri",
	},
	"mr": map[string]string{
		"name":       "marathi (marāṭhī)",
		"native":     "मराठी",
		"iso639-1":   "mr",
		"iso639-2/t": "mar",
		"iso639-2/b": "mar",
		"iso639-3":   "mar",
	},
	"mh": map[string]string{
		"name":       "marshallese",
		"native":     "kajin m̧ajeļ",
		"iso639-1":   "mh",
		"iso639-2/t": "mah",
		"iso639-2/b": "mah",
		"iso639-3":   "mah",
	},
	"mn": map[string]string{
		"name":       "mongolian",
		"native":     "монгол",
		"iso639-1":   "mn",
		"iso639-2/t": "mon",
		"iso639-2/b": "mon",
		"iso639-3":   "mon + 2",
	},
	"na": map[string]string{
		"name":       "nauru",
		"native":     "ekakairũ naoero",
		"iso639-1":   "na",
		"iso639-2/t": "nau",
		"iso639-2/b": "nau",
		"iso639-3":   "nau",
	},
	"nv": map[string]string{
		"name":       "navajo, navaho",
		"native":     "diné bizaad, dinékʼehǰí",
		"iso639-1":   "nv",
		"iso639-2/t": "nav",
		"iso639-2/b": "nav",
		"iso639-3":   "nav",
	},
	"nd": map[string]string{
		"name":       "northern ndebele",
		"native":     "isindebele",
		"iso639-1":   "nd",
		"iso639-2/t": "nde",
		"iso639-2/b": "nde",
		"iso639-3":   "nde",
	},
	"ne": map[string]string{
		"name":       "nepali",
		"native":     "नेपाली",
		"iso639-1":   "ne",
		"iso639-2/t": "nep",
		"iso639-2/b": "nep",
		"iso639-3":   "nep",
	},
	"ng": map[string]string{
		"name":       "ndonga",
		"native":     "owambo",
		"iso639-1":   "ng",
		"iso639-2/t": "ndo",
		"iso639-2/b": "ndo",
		"iso639-3":   "ndo",
	},
	"nb": map[string]string{
		"name":       "norwegian bokmål",
		"native":     "norsk bokmål",
		"iso639-1":   "nb",
		"iso639-2/t": "nob",
		"iso639-2/b": "nob",
		"iso639-3":   "nob",
	},
	"nn": map[string]string{
		"name":       "norwegian nynorsk",
		"native":     "norsk nynorsk",
		"iso639-1":   "nn",
		"iso639-2/t": "nno",
		"iso639-2/b": "nno",
		"iso639-3":   "nno",
	},
	"no": map[string]string{
		"name":       "norwegian",
		"native":     "norsk",
		"iso639-1":   "no",
		"iso639-2/t": "nor",
		"iso639-2/b": "nor",
		"iso639-3":   "nor + 2",
	},
	"ii": map[string]string{
		"name":       "nuosu",
		"native":     "ꆈꌠ꒿ nuosuhxop",
		"iso639-1":   "ii",
		"iso639-2/t": "iii",
		"iso639-2/b": "iii",
		"iso639-3":   "iii",
	},
	"nr": map[string]string{
		"name":       "southern ndebele",
		"native":     "isindebele",
		"iso639-1":   "nr",
		"iso639-2/t": "nbl",
		"iso639-2/b": "nbl",
		"iso639-3":   "nbl",
	},
	"oc": map[string]string{
		"name":       "occitan",
		"native":     "occitan, lenga d\"òc",
		"iso639-1":   "oc",
		"iso639-2/t": "oci",
		"iso639-2/b": "oci",
		"iso639-3":   "oci",
	},
	"oj": map[string]string{
		"name":       "ojibwe, ojibwa",
		"native":     "ᐊᓂᔑᓈᐯᒧᐎᓐ",
		"iso639-1":   "oj",
		"iso639-2/t": "oji",
		"iso639-2/b": "oji",
		"iso639-3":   "oji + 7",
	},
	"cu": map[string]string{
		"name":       "old church slavonic, church slavonic, old bulgarian",
		"native":     "ѩзыкъ словѣньскъ",
		"iso639-1":   "cu",
		"iso639-2/t": "chu",
		"iso639-2/b": "chu",
		"iso639-3":   "chu",
	},
	"om": map[string]string{
		"name":       "oromo",
		"native":     "afaan oromoo",
		"iso639-1":   "om",
		"iso639-2/t": "orm",
		"iso639-2/b": "orm",
		"iso639-3":   "orm + 4",
	},
	"or": map[string]string{
		"name":       "oriya",
		"native":     "ଓଡ଼ିଆ",
		"iso639-1":   "or",
		"iso639-2/t": "ori",
		"iso639-2/b": "ori",
		"iso639-3":   "ori",
	},
	"os": map[string]string{
		"name":       "ossetian, ossetic",
		"native":     "ирон æвзаг",
		"iso639-1":   "os",
		"iso639-2/t": "oss",
		"iso639-2/b": "oss",
		"iso639-3":   "oss",
	},
	"pa": map[string]string{
		"name":       "panjabi, punjabi",
		"native":     "ਪੰਜਾਬੀ, پنجابی‎",
		"iso639-1":   "pa",
		"iso639-2/t": "pan",
		"iso639-2/b": "pan",
		"iso639-3":   "pan",
	},
	"pi": map[string]string{
		"name":       "pāli",
		"native":     "पाऴि",
		"iso639-1":   "pi",
		"iso639-2/t": "pli",
		"iso639-2/b": "pli",
		"iso639-3":   "pli",
	},
	"fa": map[string]string{
		"name":       "persian (farsi)",
		"native":     "فارسی",
		"iso639-1":   "fa",
		"iso639-2/t": "fas",
		"iso639-2/b": "per",
		"iso639-3":   "fas + 2",
	},
	"pl": map[string]string{
		"name":       "polish",
		"native":     "język polski, polszczyzna",
		"iso639-1":   "pl",
		"iso639-2/t": "pol",
		"iso639-2/b": "pol",
		"iso639-3":   "pol",
	},
	"ps": map[string]string{
		"name":       "pashto, pushto",
		"native":     "پښتو",
		"iso639-1":   "ps",
		"iso639-2/t": "pus",
		"iso639-2/b": "pus",
		"iso639-3":   "pus + 3",
	},
	"pt": map[string]string{
		"name":       "portuguese",
		"native":     "português",
		"iso639-1":   "pt",
		"iso639-2/t": "por",
		"iso639-2/b": "por",
		"iso639-3":   "por",
		"sdl":        "ptb",
	},
	"qu": map[string]string{
		"name":       "quechua",
		"native":     "runa simi, kichwa",
		"iso639-1":   "qu",
		"iso639-2/t": "que",
		"iso639-2/b": "que",
		"iso639-3":   "que + 44",
	},
	"rm": map[string]string{
		"name":       "romansh",
		"native":     "rumantsch grischun",
		"iso639-1":   "rm",
		"iso639-2/t": "roh",
		"iso639-2/b": "roh",
		"iso639-3":   "roh",
	},
	"rn": map[string]string{
		"name":       "kirundi",
		"native":     "ikirundi",
		"iso639-1":   "rn",
		"iso639-2/t": "run",
		"iso639-2/b": "run",
		"iso639-3":   "run",
	},
	"ro": map[string]string{
		"name":       "romanian",
		"native":     "limba română",
		"iso639-1":   "ro",
		"iso639-2/t": "ron",
		"iso639-2/b": "rum",
		"iso639-3":   "ron",
	},
	"ru": map[string]string{
		"name":       "russian",
		"native":     "русский язык",
		"iso639-1":   "ru",
		"iso639-2/t": "rus",
		"iso639-2/b": "rus",
		"iso639-3":   "rus",
	},
	"sa": map[string]string{
		"name":       "sanskrit (saṁskṛta)",
		"native":     "संस्कृतम्",
		"iso639-1":   "sa",
		"iso639-2/t": "san",
		"iso639-2/b": "san",
		"iso639-3":   "san",
	},
	"sc": map[string]string{
		"name":       "sardinian",
		"native":     "sardu",
		"iso639-1":   "sc",
		"iso639-2/t": "srd",
		"iso639-2/b": "srd",
		"iso639-3":   "srd + 4",
	},
	"sd": map[string]string{
		"name":       "sindhi",
		"native":     "सिन्धी, سنڌي، سندھی‎",
		"iso639-1":   "sd",
		"iso639-2/t": "snd",
		"iso639-2/b": "snd",
		"iso639-3":   "snd",
	},
	"se": map[string]string{
		"name":       "northern sami",
		"native":     "davvisámegiella",
		"iso639-1":   "se",
		"iso639-2/t": "sme",
		"iso639-2/b": "sme",
		"iso639-3":   "sme",
	},
	"sm": map[string]string{
		"name":       "samoan",
		"native":     "gagana fa\"a samoa",
		"iso639-1":   "sm",
		"iso639-2/t": "smo",
		"iso639-2/b": "smo",
		"iso639-3":   "smo",
	},
	"sg": map[string]string{
		"name":       "sango",
		"native":     "yângâ tî sängö",
		"iso639-1":   "sg",
		"iso639-2/t": "sag",
		"iso639-2/b": "sag",
		"iso639-3":   "sag",
	},
	"sr": map[string]string{
		"name":       "serbian",
		"native":     "српски језик",
		"iso639-1":   "sr",
		"iso639-2/t": "srp",
		"iso639-2/b": "srp",
		"iso639-3":   "srp",
	},
	"gd": map[string]string{
		"name":       "scottish gaelic, gaelic",
		"native":     "gàidhlig",
		"iso639-1":   "gd",
		"iso639-2/t": "gla",
		"iso639-2/b": "gla",
		"iso639-3":   "gla",
	},
	"sn": map[string]string{
		"name":       "shona",
		"native":     "chishona",
		"iso639-1":   "sn",
		"iso639-2/t": "sna",
		"iso639-2/b": "sna",
		"iso639-3":   "sna",
	},
	"si": map[string]string{
		"name":       "sinhala, sinhalese",
		"native":     "සිංහල",
		"iso639-1":   "si",
		"iso639-2/t": "sin",
		"iso639-2/b": "sin",
		"iso639-3":   "sin",
	},
	"sk": map[string]string{
		"name":       "slovak",
		"native":     "slovenčina, slovenský jazyk",
		"iso639-1":   "sk",
		"iso639-2/t": "slk",
		"iso639-2/b": "slo",
		"iso639-3":   "slk",
	},
	"sl": map[string]string{
		"name":       "slovene",
		"native":     "slovenski jezik, slovenščina",
		"iso639-1":   "sl",
		"iso639-2/t": "slv",
		"iso639-2/b": "slv",
		"iso639-3":   "slv",
	},
	"so": map[string]string{
		"name":       "somali",
		"native":     "soomaaliga, af soomaali",
		"iso639-1":   "so",
		"iso639-2/t": "som",
		"iso639-2/b": "som",
		"iso639-3":   "som",
	},
	"st": map[string]string{
		"name":       "southern sotho",
		"native":     "sesotho",
		"iso639-1":   "st",
		"iso639-2/t": "sot",
		"iso639-2/b": "sot",
		"iso639-3":   "sot",
	},
	"es": map[string]string{
		"name":       "spanish, castilian",
		"native":     "español, castellano",
		"iso639-1":   "es",
		"iso639-2/t": "spa",
		"iso639-2/b": "spa",
		"iso639-3":   "spa",
	},
	"su": map[string]string{
		"name":       "sundanese",
		"native":     "basa sunda",
		"iso639-1":   "su",
		"iso639-2/t": "sun",
		"iso639-2/b": "sun",
		"iso639-3":   "sun",
	},
	"sw": map[string]string{
		"name":       "swahili",
		"native":     "kiswahili",
		"iso639-1":   "sw",
		"iso639-2/t": "swa",
		"iso639-2/b": "swa",
		"iso639-3":   "swa + 2",
	},
	"ss": map[string]string{
		"name":       "swati",
		"native":     "siswati",
		"iso639-1":   "ss",
		"iso639-2/t": "ssw",
		"iso639-2/b": "ssw",
		"iso639-3":   "ssw",
	},
	"sv": map[string]string{
		"name":       "swedish",
		"native":     "svenska",
		"iso639-1":   "sv",
		"iso639-2/t": "swe",
		"iso639-2/b": "swe",
		"iso639-3":   "swe",
	},
	"ta": map[string]string{
		"name":       "tamil",
		"native":     "தமிழ்",
		"iso639-1":   "ta",
		"iso639-2/t": "tam",
		"iso639-2/b": "tam",
		"iso639-3":   "tam",
	},
	"te": map[string]string{
		"name":       "telugu",
		"native":     "తెలుగు",
		"iso639-1":   "te",
		"iso639-2/t": "tel",
		"iso639-2/b": "tel",
		"iso639-3":   "tel",
	},
	"tg": map[string]string{
		"name":       "tajik",
		"native":     "тоҷикӣ, toğikī, تاجیکی‎",
		"iso639-1":   "tg",
		"iso639-2/t": "tgk",
		"iso639-2/b": "tgk",
		"iso639-3":   "tgk",
	},
	"th": map[string]string{
		"name":       "thai",
		"native":     "ไทย",
		"iso639-1":   "th",
		"iso639-2/t": "tha",
		"iso639-2/b": "tha",
		"iso639-3":   "tha",
	},
	"ti": map[string]string{
		"name":       "tigrinya",
		"native":     "ትግርኛ",
		"iso639-1":   "ti",
		"iso639-2/t": "tir",
		"iso639-2/b": "tir",
		"iso639-3":   "tir",
	},
	"bo": map[string]string{
		"name":       "tibetan standard, tibetan, central",
		"native":     "བོད་ཡིག",
		"iso639-1":   "bo",
		"iso639-2/t": "bod",
		"iso639-2/b": "tib",
		"iso639-3":   "bod",
	},
	"tk": map[string]string{
		"name":       "turkmen",
		"native":     "türkmen, түркмен",
		"iso639-1":   "tk",
		"iso639-2/t": "tuk",
		"iso639-2/b": "tuk",
		"iso639-3":   "tuk",
	},
	"tl": map[string]string{
		"name":       "tagalog",
		"native":     "wikang tagalog, ᜏᜒᜃᜅ᜔ ᜆᜄᜎᜓᜄ᜔",
		"iso639-1":   "tl",
		"iso639-2/t": "tgl",
		"iso639-2/b": "tgl",
		"iso639-3":   "tgl",
	},
	"tn": map[string]string{
		"name":       "tswana",
		"native":     "setswana",
		"iso639-1":   "tn",
		"iso639-2/t": "tsn",
		"iso639-2/b": "tsn",
		"iso639-3":   "tsn",
	},
	"to": map[string]string{
		"name":       "tonga", // tonga islands
		"native":     "faka tonga",
		"iso639-1":   "to",
		"iso639-2/t": "ton",
		"iso639-2/b": "ton",
		"iso639-3":   "ton",
	},
	"tr": map[string]string{
		"name":       "turkish",
		"native":     "türkçe",
		"iso639-1":   "tr",
		"iso639-2/t": "tur",
		"iso639-2/b": "tur",
		"iso639-3":   "tur",
	},
	"ts": map[string]string{
		"name":       "tsonga",
		"native":     "xitsonga",
		"iso639-1":   "ts",
		"iso639-2/t": "tso",
		"iso639-2/b": "tso",
		"iso639-3":   "tso",
	},
	"tt": map[string]string{
		"name":       "tatar",
		"native":     "татар теле, tatar tele",
		"iso639-1":   "tt",
		"iso639-2/t": "tat",
		"iso639-2/b": "tat",
		"iso639-3":   "tat",
	},
	"tw": map[string]string{
		"name":       "twi",
		"native":     "twi",
		"iso639-1":   "tw",
		"iso639-2/t": "twi",
		"iso639-2/b": "twi",
		"iso639-3":   "twi",
	},
	"ty": map[string]string{
		"name":       "tahitian",
		"native":     "reo tahiti",
		"iso639-1":   "ty",
		"iso639-2/t": "tah",
		"iso639-2/b": "tah",
		"iso639-3":   "tah",
	},
	"ug": map[string]string{
		"name":       "uyghur, uighur",
		"native":     "uyƣurqə, ئۇيغۇرچە‎",
		"iso639-1":   "ug",
		"iso639-2/t": "uig",
		"iso639-2/b": "uig",
		"iso639-3":   "uig",
	},
	"uk": map[string]string{
		"name":       "ukrainian",
		"native":     "українська мова",
		"iso639-1":   "uk",
		"iso639-2/t": "ukr",
		"iso639-2/b": "ukr",
		"iso639-3":   "ukr",
	},
	"ur": map[string]string{
		"name":       "urdu",
		"native":     "اردو",
		"iso639-1":   "ur",
		"iso639-2/t": "urd",
		"iso639-2/b": "urd",
		"iso639-3":   "urd",
	},
	"uz": map[string]string{
		"name":       "uzbek",
		"native":     "o‘zbek, ўзбек, أۇزبېك‎",
		"iso639-1":   "uz",
		"iso639-2/t": "uzb",
		"iso639-2/b": "uzb",
		"iso639-3":   "uzb + 2",
	},
	"ve": map[string]string{
		"name":       "venda",
		"native":     "tshivenḓa",
		"iso639-1":   "ve",
		"iso639-2/t": "ven",
		"iso639-2/b": "ven",
		"iso639-3":   "ven",
	},
	"vi": map[string]string{
		"name":       "vietnamese",
		"native":     "tiếng việt",
		"iso639-1":   "vi",
		"iso639-2/t": "vie",
		"iso639-2/b": "vie",
		"iso639-3":   "vie",
	},
	"vo": map[string]string{
		"name":       "volapük",
		"native":     "volapük",
		"iso639-1":   "vo",
		"iso639-2/t": "vol",
		"iso639-2/b": "vol",
		"iso639-3":   "vol",
	},
	"wa": map[string]string{
		"name":       "walloon",
		"native":     "walon",
		"iso639-1":   "wa",
		"iso639-2/t": "wln",
		"iso639-2/b": "wln",
		"iso639-3":   "wln",
	},
	"cy": map[string]string{
		"name":       "welsh",
		"native":     "cymraeg",
		"iso639-1":   "cy",
		"iso639-2/t": "cym",
		"iso639-2/b": "wel",
		"iso639-3":   "cym",
	},
	"wo": map[string]string{
		"name":       "wolof",
		"native":     "wollof",
		"iso639-1":   "wo",
		"iso639-2/t": "wol",
		"iso639-2/b": "wol",
		"iso639-3":   "wol",
	},
	"fy": map[string]string{
		"name":       "western frisian",
		"native":     "frysk",
		"iso639-1":   "fy",
		"iso639-2/t": "fry",
		"iso639-2/b": "fry",
		"iso639-3":   "fry",
	},
	"xh": map[string]string{
		"name":       "xhosa",
		"native":     "isixhosa",
		"iso639-1":   "xh",
		"iso639-2/t": "xho",
		"iso639-2/b": "xho",
		"iso639-3":   "xho",
	},
	"yi": map[string]string{
		"name":       "yiddish",
		"native":     "ייִדיש",
		"iso639-1":   "yi",
		"iso639-2/t": "yid",
		"iso639-2/b": "yid",
		"iso639-3":   "yid + 2",
	},
	"yo": map[string]string{
		"name":       "yoruba",
		"native":     "yorùbá",
		"iso639-1":   "yo",
		"iso639-2/t": "yor",
		"iso639-2/b": "yor",
		"iso639-3":   "yor",
	},
	"za": map[string]string{
		"name":       "zhuang, chuang",
		"native":     "saɯ cueŋƅ, saw cuengh",
		"iso639-1":   "za",
		"iso639-2/t": "zha",
		"iso639-2/b": "zha",
		"iso639-3":   "zha + 16",
	},
	"zu": map[string]string{
		"name":       "zulu",
		"native":     "isizulu",
		"iso639-1":   "zu",
		"iso639-2/t": "zul",
		"iso639-2/b": "zul",
		"iso639-3":   "zul",
	},
	// extra codes
	"zh-CN": map[string]string{
		"name":    "chinese simplified",
		"code":    "zh-CN",
		"bing":    "zh-CHS",
		"frengly": "zhCN",
		"promt":   "zhcn",
		"systran": "zh-Hans",
		"sdl":     "chi",
		"treu":    "zh",
	},
	"zh-TW": map[string]string{
		"name":    "chinese traditional",
		"code":    "zh-TW",
		"bing":    "zh-CHT",
		"sdl":     "cht",
		"frengly": "zhTW",
		"systran": "zh-Hant",
	},
	"haw": map[string]string{
		"name": "hawaiian",
		"code": "haw",
	},
	"iw": map[string]string{
		"name": "hebrew",
		"code": "iw",
	},
	"hmn": map[string]string{
		"name": "hmong",
		"code": "hmn",
		"bing": "mww",
	},
	"jw": map[string]string{
		"name": "javanese",
		"code": "jw",
	},
	"yue": map[string]string{
		"name": "cantonese",
		"code": "yue",
		"bing": "yue",
	},
	"tlh": map[string]string{
		"name": "klingon",
		"code": "tlh",
		"bing": "tlh",
	},
	"fil": map[string]string{
		"name": "filipino",
		"code": "fil",
		"bing": "fil",
	},
	"yua": map[string]string{
		"name": "yucateco",
		"code": "yua",
		"bing": "yua",
	},
	"otq": map[string]string{
		"name": "otomi",
		"code": "otq",
		"bing": "otq",
	},
	"mhr": map[string]string{
		"name":   "eastern mari",
		"code":   "mhr",
		"yandex": "mhr",
	},
	"mrj": map[string]string{
		"name":   "western mari",
		"code":   "mrg",
		"yandex": "mrj",
	},
	"pap": map[string]string{
		"name":   "papiamento",
		"code":   "pap",
		"yandex": "pap",
	},
	"sjn": map[string]string{
		"name":   "sindarin",
		"code":   "sjn",
		"yandex": "sjn",
	},
	"udm": map[string]string{
		"name":   "udmurt",
		"code":   "udm",
		"yandex": "udm",
	},
	"ptp": map[string]string{
		"name": "patep",
		"code": "ptp",
		"sdl":  "ptp",
	},
	"ene": map[string]string{
		"name": "enets",
		"code": "ene",
		"sdl":  "ene",
	},
	"fad": map[string]string{
		"name": "wagi",
		"code": "fad",
		"sdl":  "fad",
	},
	"ast": map[string]string{
		"name": "asturian",
		"code": "ast",
		"treu": "ast",
	},
	"dr": map[string]string{
		"name":    "???",
		"code":    "dr",
		"systran": "dr",
	},
}

type LanguageCode struct {
	Langs *map[string]map[string]string
}

func (lc *LanguageCode) Convert(value string) string {
	// Search $value
	index := ""
search:
	for key, langs := range *lc.Langs {
		for _, lang := range langs {
			if value == lang {
				index = key
				break search
			}
		}
	}

	if index == "" && value != "emj" && value != "-" { // don't care about emojis...
		log.Print("Language code " + value + " not found.")
	}

	if _, ok := (*lc.Langs)[index]["iso639-1"]; !ok {
		if _, ok := (*lc.Langs)[index]["code"]; !ok {
			return ""
		} else {
			return (*lc.Langs)[index]["code"]
		}
	}

	return (*lc.Langs)[index]["iso639-1"]
}

var Lc = &LanguageCode{&Langs}
