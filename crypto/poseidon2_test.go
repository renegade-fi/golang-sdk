package crypto

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
)

// LARGE_TEST_VECTOR is a large random vector for testing purposes
var LARGE_TEST_VECTOR = []string{
	"19230904647065636396038552256938968853322683750200788848331370785095198876878",
	"15911637536990531499916323601087363308229536839479434470684150236988586364318",
	"19564613450635237900605484014415855639638711726565938233242326595431025719394",
	"8881485276397842466310463542763687064773282714434842435007200462186454330291",
	"19360777361788036618323128391363893942287043390010641429958576638596373367481",
	"3085331459948998061513165483483441366722265041218635735749566604963071818746",
	"21558748097106774157730799228858177692912589500763760165682032362695619369317",
	"18035072001174003629301486951831250770812039363353386414795512520555897784304",
	"2661979592490645144492961866716026874439178323409329170137411752321833322685",
	"10292055339055739515033686265489930544504310589548101822739389140486942063719",
	"12994945066724721962493512176419549304629429454793867165053014941668091129119",
	"3336934521908181114356792837694897700610497660959831715955316093676622053148",
	"10516793108386028362680921580798803144530843173536722955376570406704136068089",
	"12697041080750774271917502022265858581672962400311359806762191964177668977355",
	"10094842054947678936496767847653479408619597564340716861711085790610769007265",
	"10976786782102283274879140505340730477256387596776488672647844564023117610249",
	"9845030760241303495654545376927458652242077648130407758239715702529218798316",
	"1250389595875315087131686410630559331231015553073916990281642985972353763591",
	"667445177988991638378216743507836527234666177829604929458112338327701952175",
	"6126716100898367204537998865302784056009263904334949203928563945476439860439",
	"15972857138183227620293533025780314285861761832039969919833972995379079791628",
	"20809723906700147200996803284898983115602086678608660935952851415678615140372",
	"708898091603573484211814686076051015194663342862850577190822592763033165009",
	"3630098226820926027584359965328802675147391163648927232076501047181581562514",
	"8641438173389991679614437448526630439560191300021381335279279915565862132514",
	"21118735750887348944225692755655054994945758020540257873472843820105907896813",
	"11212303057616550325894809809227038284861396363286975477316317397275331339718",
	"19239305204026369684479955281612212223414926919559298692095227051494990612093",
	"21858796315214253324561215472492723761015603710922857079718233671427249893483",
	"19541642149807315755158704273265994704492772104767432913658629111902399902844",
	"3680049808610579312375907535979426084316918928049801983698130899151101493047",
	"18056562871258487575004101084596058864047054413409356802406784098752706602420",
	"14681230599212120654407388135106787790071055592349336103290680369978147079059",
	"483726928528275754527934638201434805708674637976458305524357330445352586535",
	"14434265501379777276014526572273930799511109792700675041723257750312164570194",
	"14656851643022439568069682332210015203247773960104843202426492784839567801661",
	"13270982312799489558863384354371716760912162215598862411875585544367593796034",
	"16018698973831465607704863313952486572146279678564413108340350169970825001777",
	"14515270484327018829547624647574876047948129811997621098394469668368437300712",
	"11195191045709632014012618895426459409917495706723871649349195924635593457429",
	"3824760775693165144828621766946193728125365528199731713774823729791343435898",
	"20552846603970103813353572964325818991797822478517032075811494938286070582593",
	"10506210720662199830977403771629725636882557521994047895879001752419006069227",
	"556511866793546195714523977573308651573474098799597630174690574830712651458",
	"9027249164638657060023845280668208212129780321281899021289589873272947041509",
	"10429979893038449799921549088991016363987452892830610986118103729791010015853",
	"14642302684295266911183033128731655140462510164411933660744014724145415628032",
	"20088351040033735180924187317479868268936722467378863564857297564833890371536",
	"19464121490948166827403726618130039032375253329170596809926196142238954875645",
	"15483119362252572074827984183376616932745321270771090280053426250298508913940",
	"2623043826653855284270548179372387305295635828460358917285846625505616558823",
	"15500111015698111483331706172381392308067317348770975517477788977560465086313",
	"10478892066846514305043932946944685649032208781562141938678765613914445134626",
	"15928387331912776759516769871658745296711013567728766428611931021533854540547",
	"19262338273945176807410212969525984542954144832743261072740351012744777100740",
	"7843979522124728475403386821193285377164129064402768000612848593772414569998",
	"11502126472323222695802726966642667678801831286975468443856045694454543251489",
	"13590306640504976332565914824706860868056838388277224743623867728918058604374",
	"2022494962634765792632948135073117535883984714663985617080824369273426380299",
	"21134409127493667944326392207476286374662874636899568512123312322938249651819",
	"11966746149307589577287589821972419369674179692939583884855754874058919037143",
	"5381828135771772027614700997667075625733634880425253631671496451490497760550",
	"9054176502157935097657160935659334789082685009688032260011280456223204803671",
	"9516089283724346577112852834384416432282340653910456015568557171418816128051",
	"19599876225621554064774295479036246560498611511071591234186482627078090968291",
	"4370546288687797861164492662946026331335480722159006456627430091423481530785",
	"782152294947164802697636712154988099701588092047258221252015234606236980674",
	"20413620786518925150859997594359153562351659012407035598918227548473972330972",
	"15244540765602778922627592119165454291237626256706428838417911017958333092131",
	"21058230651592633377899154364593900248781133212012536892490810895693204299410",
	"14734437178311270369783926018451455341195993930301953203128157136445649711574",
	"16106189330975192334642709067511525078414548451774415633901644936185084929747",
	"4294000775142455095082178047353472584380965418585612350517041357537900910881",
	"2628676131866724472953292180856988567858608587078093225357654434365707430392",
	"1335948764232462849913524118644953849752503202483198226533913507984561403419",
	"19909089997757125840629924641126908203691122233209978577976255431182912026931",
	"5662667943546750188983648250916217437668802293666376445488836755053771337443",
	"15484939781036433508930902817812080192265168923821909548888794603585992484532",
	"15192699240471178683545821115666612278980852656770969676266796514107441588151",
	"18703353972860898946418109218604605152159038330158445387390218950102584426015",
	"12060695837675922083865689176673872648595767613653871747533925664424863410921",
	"14470277441995332831958756108441281236420042411013130033346993825439370082891",
	"12500674832498626633834319513485094128069710017356613465357082892313901782572",
	"12263771568343034363833639086499650996860641767409057221658011886885041501056",
	"8363090400141035183647729197868880513746136006199811126429488024476851777604",
	"10792197397153834764915362590635956684859763886779408883832504500645873202098",
	"18691575785908581259577197036575402586362601163820050541231772504632674414003",
	"7706373428731419666174073639390199324585741677535499097900463233180007890037",
	"9002588213818320709054236711814177203180322983795618797905419037909696353482",
	"17098696103362640294729882363235231267774188576939699222450902360500471028899",
	"5299655446730743145292096330695502340102208249818443375355056921935418028711",
	"3696309521969261704784993050957022486026226366671370160116286789756461496800",
	"10909447082859851945453631831476290112184878303064622229927501091474106474460",
	"4662438740728394575950653397938825496848358383385896415010943649080289721437",
	"7970993271102939000362841589122337738579289674556863446257042767568086283170",
	"1788978122459745488164484441107768254714261553830084736414481742254876427924",
	"4684834391626808342196152214182045109164528271068713733512513687680627330702",
	"17880841848614741265907462652418000551115521155627626528763013084602313801996",
	"17231236690277397900982158680966209730563708144415662739023762051723956305930",
	"19583515604391775199663263940879875683789421284095669348887925283824372088652",
}

// LARGE_TEST_VECTOR_HASH is the expected hash output for LARGE_TEST_VECTOR
var LARGE_TEST_VECTOR_HASH = "7043431630205359021101812166882265280337929418769865370612041462630759989210"

func feltFromString(s string) fr.Element {
	var felt fr.Element
	if _, err := felt.SetString(s); err != nil {
		panic(err)
	}

	return felt
}

// TestPoseidon2Sponge_Hash tests the Hash method of the Poseidon2Sponge struct.
// Test cases are derived from the Poseidon2 implementation at:
// https://github.com/renegade-fi/renegade/blob/main/renegade-crypto/src/hash/poseidon2.rs
func TestPoseidon2Sponge_Hash(t *testing.T) {
	testCases := []struct {
		name     string
		input    []fr.Element
		expected string // Expected hash result in decimal
	}{
		{
			name:     "Empty input",
			input:    []fr.Element{},
			expected: "13629302801197998987814902320299027581009939610751955228105166233386644439248",
		},
		{
			name: "Single element input 1",
			input: []fr.Element{
				fr.NewElement(1),
			},
			expected: "16195266774422401257563698575316358467855191013485223283756626417946441702527",
		},
		{
			name: "Single element randomly generated input",
			input: []fr.Element{
				feltFromString("3606381169235138002467536078418257291335960248385522353547607464602048449665"),
			},
			expected: "9974369383625325123471679280659333733739311272807606658664709120000227983436",
		},
		{
			name: "Multiple element input",
			input: []fr.Element{
				feltFromString("9047612622275400659769664160647507794672708111425750971428324879815502830296"),
				feltFromString("14016517917863009714205799673428437526088502761189873304284522050271677926200"),
				feltFromString("1141382494791791389852287594376575577791004432506599779101205944665756265553"),
				feltFromString("4392765774542063800272018853875369706251741380472033069946175621894076534574"),
				feltFromString("8620940559757973724340617232568560854674867951830226707216418953584780965297"),
				feltFromString("17221695189012834684310565357992449723137906432048741297758539210329293626569"),
				feltFromString("1717913421571771838284571386656783193322208859115951512714072329445431576879"),
				feltFromString("3016812268736145211192248591442137548205094490854033641429995837802051534264"),
				feltFromString("11626701749633546574258950888508688098221973906696870660769790107191493609035"),
				feltFromString("14645764634718958133170381507123225271145179090340545714631937853686023537646"),
			},
			expected: "9728295470246555707900495915453264867643671851345707146237532651475735824949",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sponge := NewPoseidon2Sponge()
			result := sponge.Hash(tc.input)
			assert.Equal(t, tc.expected, result.String(), "Hash result mismatch")
		})
	}
}

func TestPoseidon2Sponge_AbsorbAndSqueeze(t *testing.T) {
	sponge := NewPoseidon2Sponge()

	// Absorb some elements
	err := sponge.AbsorbBatch([]fr.Element{fr.NewElement(1), fr.NewElement(2), fr.NewElement(3)})
	assert.NoError(t, err, "AbsorbBatch should not return an error")

	// Squeeze and check results
	result1 := sponge.Squeeze()
	result2 := sponge.Squeeze()

	assert.NotEqual(t, result1, result2, "Consecutive squeezes should produce different results")
}

func TestPoseidon2Sponge_ErrorHandling(t *testing.T) {
	sponge := NewPoseidon2Sponge()

	// Absorb some elements
	err := sponge.AbsorbBatch([]fr.Element{fr.NewElement(1), fr.NewElement(2)})
	assert.NoError(t, err, "AbsorbBatch should not return an error")

	// Squeeze
	sponge.Squeeze()

	// Try to absorb after squeezing
	err = sponge.Absorb(fr.NewElement(3))
	assert.Error(t, err, "Absorb should return an error after squeezing")
}

func TestPoseidon2Sponge_LargeVector(t *testing.T) {
	// Convert the test vector to a slice of `fr.Element`
	input := make([]fr.Element, len(LARGE_TEST_VECTOR))
	for i, s := range LARGE_TEST_VECTOR {
		input[i] = feltFromString(s)
	}

	// Hash and check the result
	sponge := NewPoseidon2Sponge()
	result := sponge.Hash(input)

	assert.Equal(t, LARGE_TEST_VECTOR_HASH, result.String(), "Hash result for large test vector mismatch")
}
