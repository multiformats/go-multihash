package blake2

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/multiformats/go-multihash/core"
)

func mustHexDecode(s string) []byte {
	d, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestBlake2(t *testing.T) {
	data := mustHexDecode("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f5051525354555657")

	for _, tc := range [...]struct {
		expectedResultData   []byte
		expectedResultNoData []byte
	}{
		{mustHexDecode("57"), mustHexDecode("2e")},
		{mustHexDecode("390e"), mustHexDecode("b1fe")},
		{mustHexDecode("fe728b"), mustHexDecode("cec7ea")},
		{mustHexDecode("4fe6f30f"), mustHexDecode("1271cf25")},
		{mustHexDecode("cbc7880525"), mustHexDecode("7d64c5272e")},
		{mustHexDecode("edb9d5a35476"), mustHexDecode("ddd9c40767f9")},
		{mustHexDecode("b81be0e95423db"), mustHexDecode("4e9b03474eda9a")},
		{mustHexDecode("fa9ad8e1f5618c30"), mustHexDecode("e4a6a0577479b2b4")},
		{mustHexDecode("633306b8fb4c7a8345"), mustHexDecode("d6bd6fc9a3324e5f32")},
		{mustHexDecode("2ffc171771c4b8b9b16d"), mustHexDecode("6fa1d8fcfd719046d762")},
		{mustHexDecode("9c082c2e1484e9ec959173"), mustHexDecode("eb6ec15daf9546254f0809")},
		{mustHexDecode("b6dac9c8fd27fd6079778349"), mustHexDecode("b8e1dda3ac0aa3820ad2990b")},
		{mustHexDecode("2c9fce7a5b07fc61c9adce6d00"), mustHexDecode("50b4dc6f148a3f25b974e5c829")},
		{mustHexDecode("fbc109122fe2fa7eeafa36d6ca77"), mustHexDecode("4b1f3c22056a5cf9a3300407d264")},
		{mustHexDecode("0f69b659f66b22fbf4e34992d74c67"), mustHexDecode("b7db87196c483405e40f8401fa1fc9")},
		{mustHexDecode("907d743b004a882fc18d23df3ac0a5bb"), mustHexDecode("cae66941d9efbd404e4d88758ea67670")},
		{mustHexDecode("dabb9c4ad42c46bf780306fcd4d053d373"), mustHexDecode("246c0442cd564aced8145b8b60f1370aa7")},
		{mustHexDecode("d7c79aaa168ec2c806756822486a3b31786c"), mustHexDecode("91a1a481a82eb3f3e6262de11f142d234945")},
		{mustHexDecode("0cec62c268f2b747b050f29e8a4713c4dde635"), mustHexDecode("35bd4214446fda5ce2e05015f1ba43e26f1b96")},
		{mustHexDecode("c5b9110bfb2da0ce0a55ba3f2be0aaab233c5a14"), mustHexDecode("3345524abf6bbe1809449224b5972c41790b6cf2")},
		{mustHexDecode("9eb92b042d3c1f2a03441c0bafe53d93991c1466bf"), mustHexDecode("077d8272052a6edfff4047461c3a2b3d9d330dbbf0")},
		{mustHexDecode("7043e1b53d58f38a1262382a0f3b8c13289752bba647"), mustHexDecode("1065c75a5ab372acff0b521808a4766c70b12b10ad8c")},
		{mustHexDecode("9c9f549599752e7f4227c84f77d2156d373d57a7b31cd4"), mustHexDecode("e30b37bb45ad2f1954a0ab31666f909df8d4eabd6933e9")},
		{mustHexDecode("bbdd0013817106ca72be7f23665acc5d8ac5e7cbe9b8316d"), mustHexDecode("ab3b5331a7135ed50d0f182d026e60abdb3646fd51bcf8a3")},
		{mustHexDecode("0f928ecf8557ddd62b0f5fe9252632c7dd7e88e6f4faf5edf2"), mustHexDecode("94165bbe7a8a0f49fad8c1b39c40b7dd613409378dcc47681f")},
		{mustHexDecode("93189ce79eb1257ecfc01a6d2d0fa345beaa34610af9ee9c1240"), mustHexDecode("7895f50fee886d460f321601da8d2db483a08c0264cd8ff3617e")},
		{mustHexDecode("a6ff723c4e43a3b46655152918422485bdf8327637cf1078476d10"), mustHexDecode("b41793f77a58236ee36d36570bcd14cf00ba6a443c6c5bd4bb9eaf")},
		{mustHexDecode("2e30811d4d0942f5e2bd6f52961efe6f0e583ba39ddd368e004200ab"), mustHexDecode("836cc68931c2e4e3e838602eca1902591d216837bafddfe6f0c8cb07")},
		{mustHexDecode("ed641f9462d220d0aebd649fcdf0af533e4cd5a9762645c226e56e4659"), mustHexDecode("a10eae68c06d70c597699d656d6ae213430569f9c62e04cd2fc3a0c1bf")},
		{mustHexDecode("e027addc1bdf2d24f8f763a22906ac457d199f8c3c3f91c20cfad7e20e55"), mustHexDecode("a5d6d5975d09c76462b3f9c74f9568d9f9fd46dfbdcbf3f14bc835298b22")},
		{mustHexDecode("5640281fd6fd8302a7773814289f5170790df63f263d9a8c6c642c8c705514"), mustHexDecode("b4d7d8f500d546e71fe03f080b6bfefd567a0aa97e84bdb2cf8b15d1867c00")},
		{mustHexDecode("fd9664ecb814785c8474188706e6ab0952925f9b9d8e351665ced12e84d92fad"), mustHexDecode("0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8")},
		{mustHexDecode("972b75e4a5747c695af8c430e685d2090c766a067ac629f5d75ccaa532de6aab24"), mustHexDecode("ddca500c4d28f7f2816de1574f840e4878c1c5aa30c149745e0149273b214c359d")},
		{mustHexDecode("6abf92bc9a9747942157afb97af768b798c1c2dd54188db396c529e6af07fcbba5bd"), mustHexDecode("90933ab63c7665e2bd6431e496ec60d38839fbec78e33aae2c152c073f64264bdab9")},
		{mustHexDecode("2b5c40e57ec34bc2d1a5759e46870d7324fdc0f55208492b31d637dc42e2de8a26ce18"), mustHexDecode("148833bb2bfcc18b9e90024eaeecc0a96027a777761e0b9c93d6642937bb4b8705e218")},
		{mustHexDecode("db0c682edad98e563aa47646050d812f0c0e72c21e52eb4f089cdc48fe8e0c649ec403f5"), mustHexDecode("92f3592c601fe36aa32c62e305f965905a2982dee6a45c09011ddf05f9cf9b7b5609414f")},
		{mustHexDecode("0f78215590aff9081907909e4fb3ed7ee0aff65415006aa1ff5d45b670a05b1cba03afaba6"), mustHexDecode("6d82c523a958c2b00e42701be980963438d5f40572c70d3d723c03ddebdb74575866f3adbb")},
		{mustHexDecode("87eba8d23f0e7104e728b9e718e622149f76ee964dff4d752d7f1a4370e7f8d0f04d8bcd5787"), mustHexDecode("dc5abbc8c533139ba5873c9562868914e501b13aadc59c143d1bfe97cbcb5fab5b65ed488158")},
		{mustHexDecode("1c4481143ad1fa3d7bba1967b74e02a6d888d948573b7aef9d0b9501513366bf4c7e642f65ce36"), mustHexDecode("61a54c550005791e4726043fbfc347bb8952e520818157aeaf0d0f877c51950e06ff3157d02a6f")},
		{mustHexDecode("7c6d4ba464a293c8e51402abb64ab02c0792d4d961769fbe0dce5134e20ec28d0e48ed0204dfedb4"), mustHexDecode("2e316d2c76c9760df1e604e4ffd1aa5ac6c6ac50aaa8071f7313ea931e205da084bbae9a2019f6aa")},
		{mustHexDecode("1089aa00e244d33468840ba902b93f723e969b969afe728e79485c7eccf2a9ec3105dffa625ffee3ec"), mustHexDecode("592c90e91f3187c352649476b86bba76c128433e6f3ac8c75710042f4b310e1c7aea39b0aff9b51bd3")},
		{mustHexDecode("2fcd69fb4270e081a222aab14a22315a7d862798000f548082b0fe02d01e37218edd6305ccbea1a35626"), mustHexDecode("f564703984efb278dfb04536d0bf4b86a17e8a9847104f773b81835ffc60b343a364e224e36552728dd6")},
		{mustHexDecode("6a5c37e5490fa3988894104ce266c951aca5002cd27f4ab8aa31eed60d66ee9fbb016d71fff79d4290fd24"), mustHexDecode("5112353efd2617941caf7de611f152ac7b6fbacfb682aa43ecb707c8977ae8f307e50da1942c6eed777082")},
		{mustHexDecode("2cdbbad7908e4cd48251aff565dd0cc5a17deb061c10343b2f539b44978982d3b9e7289cb167aebf48c53496"), mustHexDecode("b2e01f2639b7e74abab0bb7e88f7ab7ae94ba6292c3a42537ca288635259a50edd9c7d7a1c7b8d2e2f86848e")},
		{mustHexDecode("590a0f21f0f6019bdfb1e53b76caacc8131ed402433ba425b74aacfd9858ac9fe6df644b293509392ef305d7ee"), mustHexDecode("fa9d9e37d6fe09eb8116510fadb9c61cc59e332d46cc4a365e72edc733188f08be9c0894b6dbb06023ff312506")},
		{mustHexDecode("043211ff429d42bf06722329f2e71bbf3bc9fa6f7742b8bcc6e23cc0b77c708ed83772197fbe272c0323ced6d5bd"), mustHexDecode("d47deb78c6d8db06e3b38d8faa368d22cbab03cbfb2b3ad201be5729ab454278007f76dcdb14de4eb38958745f77")},
		{mustHexDecode("13f0a7e600fdbd814b5d3c7a19a3128ed460e23d71a66b2875eaa6bbefca778fd5b8258ac9d5f3cbe268a88798eec7"), mustHexDecode("e4ac268b5be19d515b8ddd90bc7e89100f875fa994517409907cb6f3c6eefacc3890c84dd3e91cd2886eb57033c749")},
		{mustHexDecode("5389909a11359d68263875199b2d53bc2c473d44bd5eccc562634ddaed80c6075a35bda179aa3129eaec1a11b87af6a8"), mustHexDecode("b32811423377f52d7862286ee1a72ee540524380fda1724a6f25d7978c6fd3244a6caf0498812673c5e05ef583825100")},
		{mustHexDecode("d4815007ab460e95a75e2f4912b8b6bcaeefb18e81b473a58663a85da25e680cb22d965b2dacb7cd69b294b77c3a85c74c"), mustHexDecode("a993b7c6dbd66f7a45487707d7e3eda19201f7fec9dcf1ae3c0a66eb4be4d21ed8af10490cef0c3168e9ff0dcfb5dcd651")},
		{mustHexDecode("b0828379fa986f452816cf9e98cbd3b7b66e7130512db924c123c465898761e8ae390b94727a782834df4f77ada7d49f9508"), mustHexDecode("3189e5764c09a2f5d1d9f5cb1967ebd3dfeade9c62af8bb0dc032bb3e90dd1e760fbaba8956f97c7602d0a2ec162169ef219")},
		{mustHexDecode("c09732d4ab7631a57e3a7b94fd50bb281318b4a7eb8903b7179caec2e7d0a8bc507ccb0f14593c69318ccf2d39bc3073e564cd"), mustHexDecode("31635ed8064b99e056ed7009905673c986944a718c6e5935e7eeb67652550d56fe7ec110a383ef94ef7977be456a44503434ad")},
		{mustHexDecode("22951e6c2fadbb2e9f3615895913a8e5a034d1e5dbfb7c168c2a89121be56e6ee94243b94e672452050bae6e3b53177398e0fac2"), mustHexDecode("f4e2de2be49787b13e0b38c0d02578b78a76f6c8fc48948c00f67812bd6c9ceaff17b04617532862be3cb251524b93d83a266e35")},
		{mustHexDecode("35ab960da50b0d2bc0ff05a59a10fcfd58b1484c87251f872dea2a9957304e9cea84b7ec07e9eb2f716a3d16bd97693d37a79d1253"), mustHexDecode("e3af5d079bce8fbbad6f5047d77025b8e100d91ecc066fa525d290ef6a867f93b2798769067f8790df954682011617a68d7169ef15")},
		{mustHexDecode("cfd083d366fab80952cfbcf2988edd1c15109870cf3296f0f9e04d4ca50610ef4dfc73a5cd28586fecc2e3e4d7c60c07b9621dcf877d"), mustHexDecode("0668149330f455fe58c70d209cff452742cc1125eee5e1d67af18e9b2a67b5ca6973940135341c2807c9237295ec0a0d173dbc28f687")},
		{mustHexDecode("e2fc9688e8fa0f63dd4038157d7476df5c0a6bdba9fb95d971646f2146adf5b9e4299306c38660fec59a392d503600e62c6abf99df92ec"), mustHexDecode("89c4f154fddb635864729c086c40ff2e574ef4fa1ab592d9bee584693852cfeee57c743b9a8771443e522f454218b260838c0a913d29e5")},
		{mustHexDecode("232911189422198fa470df67ba475afb61590c168e8f89744ad69e959059e168186b5344b86dd0b39f715e9ca7e922d23b0a5c29b4287962"), mustHexDecode("e7d2cb731e704ab61a3fa0ddd3bb3a6bfe3c3bc03b2c80a7545a0c9cedb575dfaa6821be9879e9ecd24350297f14470ad3d1cd2d19f27fbf")},
		{mustHexDecode("e0ab326421e4732381ab2ef306d0e3751ecb537d2ec2020f87f3e15d06a263239b824d7f5921b13d03c9b383191cb5fd9f9bda879264b730a3"), mustHexDecode("a6e2604d330fa35f9f97cb89a4160928704e058f1aa0badc51b6e16afa943362fc1b32a4d79138b8103dfcad3239de59c17a267e72f7a0693e")},
		{mustHexDecode("87de3953f266e1d4bddb48d0c416daeb83f1b6631c79db5886856749dc3da4e7e541ad2510c5c3332838fafe98921d6005851d22773ecff551db"), mustHexDecode("f3cc91641a39f6acada71544227505ae109b8c86c2f5fc3c4b7265c64ca6e99967824cea78f6ffb9a0851c86aa52b28ba3352164eedfefc80ddd")},
		{mustHexDecode("1f01b65d74228e0ce89c98f43fb0c39868f8be8afbceb264cfc1adafebd58eae74d666fb184af4e69cac74b1c7561ab172b0ff8024a344e148e901"), mustHexDecode("cf1335ff92a6710c3cfa3dd8ac8c7a435aece775997bdaa1ac57276b0fa16b9a5f1f78a334eefafd0bc9d9cafa6633ba7abed8f67ce8d287af1822")},
		{mustHexDecode("08f648dbc6b47dae03dd4b537926b02e7963679e096b2276be5200769484ddc9dd02ad8d342841dcb29916d8964ede7f09e2ecd6deb3910e5e2e36d1"), mustHexDecode("22f194f655ea58d7fefe35b09c91c91cf5e1a4047181ea7cd7674e597be65f6541fa1fdddf404e7851b1d471478048d550546d14d88345fb422c19f6")},
		{mustHexDecode("dba08fdd0908deb20de12f7c54dc139457c2d6d623145ba2f6f7e01e7d3d541bed2318378eebe337ceebea5738ba52c79a454a4fa253907fc747e3b0ed"), mustHexDecode("d10c86444347b9bbb839717bc3161a10412c52fd2eb52c0a08fcd4c1f091801c0b2b09c74d716f4874761ec1b11afd66be0e13b129b6bc877720f2c7fd")},
		{mustHexDecode("eafa5e0a0d3a9198276ed0d189d71424caf450dced14b223e6caa4d8d5667ecdf8268431f0b95fab70e7c0e8741107d03d3bc982a2b6765668799f0745f8"), mustHexDecode("50e5578cdbe722b76b9b7d629aec8fb4926b4073da62774e64cafa1b33627c24d70009660e784558b3daa7a65b6841976c41cf3d6891ea1ccdd10894e64d")},
		{mustHexDecode("ca645dcaed18a7bd232931a5e5d20d51ebc3fd667d5c7549fe8cef2b52fb8bb83546ab55243c5c1e58a07410a89d548fea5b1221971f87448a0d594e26688c"), mustHexDecode("4ded8c5fc8b12f3273f877ca585a44ad6503249a2b345d6d9c0e67d85bcb700db4178c0303e93b8f4ad758b8e2c9fd8b3d0c28e585f1928334bb77d36782e8")},
		{mustHexDecode("4665cef8ba4db4d0acb118f2987f0bb09f8f86aa445aa3d5fc9a8b346864787489e8fcecc125d17e9b56e12988eac5ecc7286883db0661b8ff05da2afff30fe4"), mustHexDecode("786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce")},
	} {
		size := uint64(len(tc.expectedResultData))
		t.Run(fmt.Sprintf("blake2b_%d", size*8), func(t *testing.T) {
			h, err := multihash.GetHasher(blake2b_min + size - 1)
			if err != nil {
				t.Errorf("failed to get: %s", err.Error())
				return
			}

			if result := h.Sum(nil); !bytes.Equal(result, tc.expectedResultNoData) {
				t.Errorf("digest empty doesn't match, expected %s; got %s", hex.EncodeToString(tc.expectedResultNoData), hex.EncodeToString(result))
			}
			n, err := h.Write(data)
			if err != nil {
				t.Errorf("hashing data failed: %s", err.Error())
				return
			}
			if n != len(data) {
				t.Errorf("not enough bytes hashed, expected %d; got %d", len(data), n)
				return
			}
			if result := h.Sum(nil); !bytes.Equal(result, tc.expectedResultData) {
				t.Errorf("digest full doesn't match, expected %s; got %s", hex.EncodeToString(tc.expectedResultData), hex.EncodeToString(result))
			}
		})
	}
}
