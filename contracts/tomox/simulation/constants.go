package simulation

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	BaseTOMO    = big.NewInt(0).Mul(big.NewInt(10), big.NewInt(100000000000000000)) // 1 TOMO
	RpcEndpoint = "http://127.0.0.1:8501/"
	MainKey, _  = crypto.HexToECDSA(os.Getenv("MAIN_ADDRESS_KEY"))
	MainAddr    = crypto.PubkeyToAddress(MainKey.PublicKey) //0x17F2beD710ba50Ed27aEa52fc4bD7Bda5ED4a037

	// TRC21 Token
	MinTRC21Apply = big.NewInt(0).Mul(big.NewInt(10), BaseTOMO) // 10 TOMO
	TRC21TokenCap = big.NewInt(0).Mul(big.NewInt(1000000000000), BaseTOMO)
	TRC21TokenFee = big.NewInt(100)

	// TOMOX
	MaxRelayers  = big.NewInt(200)
	MaxTokenList = big.NewInt(200)
	MinDeposit   = big.NewInt(0).Mul(big.NewInt(25000), BaseTOMO) // 25000 TOMO
	TradeFee     = uint16(1)

	RelayerCoinbaseKey, _ = crypto.HexToECDSA(os.Getenv("RELAYER_COINBASE_KEY")) //
	RelayerCoinbaseAddr   = crypto.PubkeyToAddress(RelayerCoinbaseKey.PublicKey) // 0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e

	OwnerRelayerKey, _ = crypto.HexToECDSA(os.Getenv("RELAYER_OWNER_KEY"))
	OwnerRelayerAddr   = crypto.PubkeyToAddress(OwnerRelayerKey.PublicKey) //0x703c4b2bD70c169f5717101CaeE543299Fc946C7

	TOMONative = common.HexToAddress("0x0000000000000000000000000000000000000001")

	TokenNameList    = []string{"BTC", "ETH", "XRP", "LTC", "BNB", "ADA", "ETC", "BCH", "EOS"}
	TokenNameAddress = map[string]common.Address{"BTC": common.HexToAddress("0x4d7eA2cE949216D6b120f3AA10164173615A2b6C"),
		"ETH": common.HexToAddress("0xC2fa1BA90b15E3612E0067A0020192938784D9C5"),
		"XRP": common.HexToAddress("0xAad540ac542C3688652a3fc7b8e21B3fC1D097e9"),
		"LTC": common.HexToAddress("0x5dc27D59bB80E0EF853Bb2e27B94113DF08F547F"),
		"BNB": common.HexToAddress("0x6F98655A8fa7AEEF3147ee002c666d09c7AA4F5c"),
		"ADA": common.HexToAddress("0xaC389aCA56394a5B14918cF6437600760B6c650C"),
		"ETC": common.HexToAddress("0x576201Ac3f1E0fe483a9320DaCc4B08EB3E58306"),
		"BCH": common.HexToAddress("0xf992cf45394dAc5f50A26446de17803a79B940da"),
		"EOS": common.HexToAddress("0xFDF68dE6dFFd893221fc9f7985FeBC2AB20761A6"),
		"TOMO":common.HexToAddress(common.TomoNativeAddress),
	}
	TeamAddresses = []common.Address{
		common.HexToAddress("0x8fB1047e874d2e472cd08980FF8383053dd83102"), // MM team
		common.HexToAddress("0x9ca1514E3Dc4059C29a1608AE3a3E3fd35900888"), // MM team
		common.HexToAddress("0x15e08dE16f534c890828F2a0D935433aF5B3CE0C"), // CTO
		common.HexToAddress("0xb68D825655F2fE14C32558cDf950b45beF18D218"), // DEX team
		common.HexToAddress("0xF7349C253FF7747Df661296E0859c44e974fb52E"), // HaiDV
		common.HexToAddress("0x9f6b8fDD3733B099A91B6D70CDC7963ebBbd2684"), // Can
		common.HexToAddress("0x06605B28aab9835be75ca242a8aE58f2e15F2F45"), // Nien
		common.HexToAddress("0x33c2E732ae7dce8B05F37B2ba0CFe14c980c4Dbe"), // Vu Pham
		common.HexToAddress("0x16a73f3a64eca79e117258e66dfd7071cc8312a9"), // BTCTOMO
		common.HexToAddress("0xac177441ac2237b2f79ecff1b8f6bca39e27ef9f"), // ETHTOMO
		common.HexToAddress("0x4215250e55984c75bbce8ae639b86a6cad8ec126"), // XRPTOMO
		common.HexToAddress("0x6b70ca959814866dd5c426d63d47dde9cc6c32d2"), // LTCTOMO
		common.HexToAddress("0x33df079fe9b9cd7fb23a1085e4eaaa8eb6952cb3"), // BNBTOMO
		common.HexToAddress("0x3cab8292137804688714670640d19f9d7a60c472"), // ADATOMO
		common.HexToAddress("0x9415d953d47c5f155cac9de7b24a756f352eafbf"), // ETCTOMO
		common.HexToAddress("0xe32d2e7c8e8809e45c8e2332830b48d9e231e3f2"), // BCHTOMO
		common.HexToAddress("0xf76ddbda664ea47088937e1cf9ff15036714dee3"), // EOSTOMO
		common.HexToAddress("0xc465ee82440dada9509feb235c7cd7d896acf13c"), // ETHBTC
		common.HexToAddress("0xb95bdc136c579dc3fd2b2424a8e925a90228d2c2"), // XRPBTC
	}
)
