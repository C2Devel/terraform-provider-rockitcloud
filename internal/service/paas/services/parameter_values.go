package services

// These values are used as defaults for the fields that can't hold 0 as default value.
// It allows to prevent sending the field until the user sets it in config.
const (
	mySQLGcsFcFactorDefault                        = -1
	mySQLInnodbThreadConcurrencyDefault            = -1
	mySQLThreadCacheSizeDefault                    = -1
	postgreSQLMaxParallelMaintenanceWorkersDefault = -1
	postgreSQLWalKeepSegmentsDefault               = -1
)

func mySQLDatabaseUserPrivileges() []string {
	return []string{
		"ALL",
		"ALTER",
		"ALTER ROUTINE",
		"CREATE",
		"CREATE ROUTINE",
		"CREATE TEMPORARY TABLES",
		"CREATE VIEW",
		"DELETE",
		"DROP",
		"EVENT",
		"EXECUTE",
		"INDEX",
		"INSERT",
		"LOCK TABLES",
		"SELECT",
		"SHOW VIEW",
		"TRIGGER",
		"UPDATE",
	}
}

func postgreSQLDatabaseExtensions() []string {
	return []string{
		"address_standardizer",
		"address_standardizer_data_us",
		"amcheck",
		"autoinc",
		"bloom",
		"btree_gin",
		"btree_gist",
		"citext",
		"cube",
		"dblink",
		"dict_int",
		"dict_xsyn",
		"earthdistance",
		"fuzzystrmatch",
		"hstore",
		"intarray",
		"isn",
		"lo",
		"ltree",
		"moddatetime",
		"pg_buffercache",
		"pg_trgm",
		"pg_visibility ",
		"pgcrypto",
		"pgrowlocks",
		"pgstattuple",
		"postgis",
		"postgis_tiger_geocoder",
		"postgis_topology",
		"postgres_fdw",
		"seg",
		"tablefunc",
		"tcn",
		"timescaledb",
		"tsm_system_rows",
		"tsm_system_time",
		"unaccent",
		"uuid-ossp",
		"xml2",
	}
}

func postgreSQLDatabaseLocales() []string {
	return []string{
		"C",
		"aa_DJ.UTF-8",
		"aa_ER.UTF-8",
		"aa_ET.UTF-8",
		"af_ZA.UTF-8",
		"am_ET.UTF-8",
		"an_ES.UTF-8",
		"ar_AE.UTF-8",
		"ar_BH.UTF-8",
		"ar_DZ.UTF-8",
		"ar_EG.UTF-8",
		"ar_IN.UTF-8",
		"ar_IQ.UTF-8",
		"ar_JO.UTF-8",
		"ar_KW.UTF-8",
		"ar_LB.UTF-8",
		"ar_LY.UTF-8",
		"ar_MA.UTF-8",
		"ar_OM.UTF-8",
		"ar_QA.UTF-8",
		"ar_SA.UTF-8",
		"ar_SD.UTF-8",
		"ar_SY.UTF-8",
		"ar_TN.UTF-8",
		"ar_YE.UTF-8",
		"as_IN.UTF-8",
		"ast_ES.UTF-8",
		"az_AZ.UTF-8",
		"be_BY.UTF-8",
		"bem_ZM.UTF-8",
		"ber_DZ.UTF-8",
		"ber_MA.UTF-8",
		"bg_BG.UTF-8",
		"bho_IN.UTF-8",
		"bn_BD.UTF-8",
		"bn_IN.UTF-8",
		"bo_CN.UTF-8",
		"bo_IN.UTF-8",
		"br_FR.UTF-8",
		"brx_IN.UTF-8",
		"bs_BA.UTF-8",
		"byn_ER.UTF-8",
		"ca_AD.UTF-8",
		"ca_ES.UTF-8",
		"ca_FR.UTF-8",
		"ca_IT.UTF-8",
		"crh_UA.UTF-8",
		"csb_PL.UTF-8",
		"cs_CZ.UTF-8",
		"cv_RU.UTF-8",
		"cy_GB.UTF-8",
		"da_DK.UTF-8",
		"de_AT.UTF-8",
		"de_BE.UTF-8",
		"de_CH.UTF-8",
		"de_DE.UTF-8",
		"de_LU.UTF-8",
		"dv_MV.UTF-8",
		"dz_BT.UTF-8",
		"el_CY.UTF-8",
		"el_GR.UTF-8",
		"en_AG.UTF-8",
		"en_AU.UTF-8",
		"en_BW.UTF-8",
		"en_CA.UTF-8",
		"en_DK.UTF-8",
		"en_GB.UTF-8",
		"en_HK.UTF-8",
		"en_IE.UTF-8",
		"en_IN.UTF-8",
		"en_NG.UTF-8",
		"en_NZ.UTF-8",
		"en_PH.UTF-8",
		"en_SG.UTF-8",
		"en_US.UTF-8",
		"en_ZA.UTF-8",
		"en_ZM.UTF-8",
		"en_ZW.UTF-8",
		"es_AR.UTF-8",
		"es_BO.UTF-8",
		"es_CL.UTF-8",
		"es_CO.UTF-8",
		"es_CR.UTF-8",
		"es_CU.UTF-8",
		"es_DO.UTF-8",
		"es_EC.UTF-8",
		"es_ES.UTF-8",
		"es_GT.UTF-8",
		"es_HN.UTF-8",
		"es_MX.UTF-8",
		"es_NI.UTF-8",
		"es_PA.UTF-8",
		"es_PE.UTF-8",
		"es_PR.UTF-8",
		"es_PY.UTF-8",
		"es_SV.UTF-8",
		"es_US.UTF-8",
		"es_UY.UTF-8",
		"es_VE.UTF-8",
		"et_EE.UTF-8",
		"eu_ES.UTF-8",
		"fa_IR.UTF-8",
		"ff_SN.UTF-8",
		"fi_FI.UTF-8",
		"fil_PH.UTF-8",
		"fo_FO.UTF-8",
		"fr_BE.UTF-8",
		"fr_CA.UTF-8",
		"fr_CH.UTF-8",
		"fr_FR.UTF-8",
		"fr_LU.UTF-8",
		"fur_IT.UTF-8",
		"fy_DE.UTF-8",
		"fy_NL.UTF-8",
		"ga_IE.UTF-8",
		"gd_GB.UTF-8",
		"gez_ER.UTF-8",
		"gez_ET.UTF-8",
		"gl_ES.UTF-8",
		"gu_IN.UTF-8",
		"gv_GB.UTF-8",
		"ha_NG.UTF-8",
		"he_IL.UTF-8",
		"hi_IN.UTF-8",
		"hne_IN.UTF-8",
		"hr_HR.UTF-8",
		"hsb_DE.UTF-8",
		"ht_HT.UTF-8",
		"hu_HU.UTF-8",
		"hy_AM.UTF-8",
		"id_ID.UTF-8",
		"ig_NG.UTF-8",
		"ik_CA.UTF-8",
		"is_IS.UTF-8",
		"it_CH.UTF-8",
		"it_IT.UTF-8",
		"iu_CA.UTF-8",
		"iw_IL.UTF-8",
		"ja_JP.UTF-8",
		"ka_GE.UTF-8",
		"kk_KZ.UTF-8",
		"kl_GL.UTF-8",
		"km_KH.UTF-8",
		"kn_IN.UTF-8",
		"kok_IN.UTF-8",
		"ko_KR.UTF-8",
		"ks_IN.UTF-8",
		"ku_TR.UTF-8",
		"kw_GB.UTF-8",
		"ky_KG.UTF-8",
		"lb_LU.UTF-8",
		"lg_UG.UTF-8",
		"li_BE.UTF-8",
		"lij_IT.UTF-8",
		"li_NL.UTF-8",
		"lo_LA.UTF-8",
		"lt_LT.UTF-8",
		"lv_LV.UTF-8",
		"mai_IN.UTF-8",
		"mg_MG.UTF-8",
		"mhr_RU.UTF-8",
		"mi_NZ.UTF-8",
		"mk_MK.UTF-8",
		"ml_IN.UTF-8",
		"mn_MN.UTF-8",
		"mr_IN.UTF-8",
		"ms_MY.UTF-8",
		"mt_MT.UTF-8",
		"my_MM.UTF-8",
		"nb_NO.UTF-8",
		"nds_DE.UTF-8",
		"nds_NL.UTF-8",
		"ne_NP.UTF-8",
		"nl_AW.UTF-8",
		"nl_BE.UTF-8",
		"nl_NL.UTF-8",
		"nn_NO.UTF-8",
		"nr_ZA.UTF-8",
		"nso_ZA.UTF-8",
		"oc_FR.UTF-8",
		"om_ET.UTF-8",
		"om_KE.UTF-8",
		"or_IN.UTF-8",
		"os_RU.UTF-8",
		"pa_IN.UTF-8",
		"pap_AN.UTF-8",
		"pa_PK.UTF-8",
		"pl_PL.UTF-8",
		"ps_AF.UTF-8",
		"pt_BR.UTF-8",
		"pt_PT.UTF-8",
		"ro_RO.UTF-8",
		"ru_RU.UTF-8",
		"ru_UA.UTF-8",
		"rw_RW.UTF-8",
		"sa_IN.UTF-8",
		"sc_IT.UTF-8",
		"sd_IN.UTF-8",
		"se_NO.UTF-8",
		"shs_CA.UTF-8",
		"sid_ET.UTF-8",
		"si_LK.UTF-8",
		"sk_SK.UTF-8",
		"sl_SI.UTF-8",
		"so_DJ.UTF-8",
		"so_ET.UTF-8",
		"so_KE.UTF-8",
		"so_SO.UTF-8",
		"sq_AL.UTF-8",
		"sq_MK.UTF-8",
		"sr_ME.UTF-8",
		"sr_RS.UTF-8",
		"ss_ZA.UTF-8",
		"st_ZA.UTF-8",
		"sv_FI.UTF-8",
		"sv_SE.UTF-8",
		"sw_KE.UTF-8",
		"sw_TZ.UTF-8",
		"ta_IN.UTF-8",
		"ta_LK.UTF-8",
		"te_IN.UTF-8",
		"tg_TJ.UTF-8",
		"th_TH.UTF-8",
		"ti_ER.UTF-8",
		"ti_ET.UTF-8",
		"tig_ER.UTF-8",
		"tk_TM.UTF-8",
		"tl_PH.UTF-8",
		"tn_ZA.UTF-8",
		"tr_CY.UTF-8",
		"tr_TR.UTF-8",
		"ts_ZA.UTF-8",
		"tt_RU.UTF-8",
		"ug_CN.UTF-8",
		"uk_UA.UTF-8",
		"unm_US.UTF-8",
		"ur_IN.UTF-8",
		"ur_PK.UTF-8",
		"ve_ZA.UTF-8",
		"vi_VN.UTF-8",
		"wa_BE.UTF-8",
		"wae_CH.UTF-8",
		"wal_ET.UTF-8",
		"wo_SN.UTF-8",
		"xh_ZA.UTF-8",
		"yi_US.UTF-8",
		"yo_NG.UTF-8",
		"yue_HK.UTF-8",
		"zh_CN.UTF-8",
		"zh_HK.UTF-8",
		"zh_SG.UTF-8",
		"zh_TW.UTF-8",
		"zu_ZA.UTF-8",
	}
}

//nolint:unused // Function will be used in validation for `encoding` field.
func postgreSQLDatabaseEncodingMap() map[string][]string {
	return map[string][]string{
		"C": {
			"EUC_CN",
			"EUC_JP",
			"EUC_JIS_2004",
			"EUC_KR",
			"EUC_TW",
			"ISO_8859_5",
			"ISO_8859_6",
			"ISO_8859_7",
			"ISO_8859_8",
			"KOI8R",
			"KOI8U",
			"LATIN1",
			"LATIN2",
			"LATIN3",
			"LATIN4",
			"LATIN5",
			"LATIN6",
			"LATIN7",
			"LATIN8",
			"LATIN9",
			"LATIN10",
			"MULE_INTERNAL",
			"SQL_ASCII",
			"UTF8",
			"WIN866",
			"WIN874",
			"WIN1250",
			"WIN1251",
			"WIN1252",
			"WIN1253",
			"WIN1254",
			"WIN1255",
			"WIN1256",
			"WIN1257",
			"WIN1258",
		},
		"aa_DJ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"aa_ER.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"aa_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"af_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"am_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"an_ES.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_AE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_BH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_DZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_EG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_IQ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_JO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_KW.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_LB.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_LY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_MA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_OM.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_QA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_SA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_SD.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_SY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_TN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ar_YE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"as_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ast_ES.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"az_AZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"be_BY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bem_ZM.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ber_DZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ber_MA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bg_BG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bho_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bn_BD.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bn_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bo_CN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bo_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"br_FR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"brx_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"bs_BA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"byn_ER.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ca_AD.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ca_ES.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ca_FR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ca_IT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"crh_UA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"csb_PL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"cs_CZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"cv_RU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"cy_GB.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"da_DK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"de_AT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"de_BE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"de_CH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"de_DE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"de_LU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"dv_MV.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"dz_BT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"el_CY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"el_GR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_AG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_AU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_BW.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_CA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_DK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_GB.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_HK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_IE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_NG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_NZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_PH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_SG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_US.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_ZM.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"en_ZW.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_AR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_BO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_CL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_CO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_CR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_CU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_DO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_EC.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_ES.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_GT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_HN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_MX.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_NI.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_PA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_PE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_PR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_PY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_SV.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_US.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_UY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"es_VE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"et_EE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"eu_ES.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fa_IR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ff_SN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fi_FI.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fil_PH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fo_FO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fr_BE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fr_CA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fr_CH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fr_FR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fr_LU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fur_IT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fy_DE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"fy_NL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ga_IE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"gd_GB.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"gez_ER.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"gez_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"gl_ES.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"gu_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"gv_GB.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ha_NG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"he_IL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"hi_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"hne_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"hr_HR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"hsb_DE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ht_HT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"hu_HU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"hy_AM.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"id_ID.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ig_NG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ik_CA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"is_IS.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"it_CH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"it_IT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"iu_CA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"iw_IL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ja_JP.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ka_GE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"kk_KZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"kl_GL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"km_KH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"kn_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"kok_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ko_KR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ks_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ku_TR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"kw_GB.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ky_KG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"lb_LU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"lg_UG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"li_BE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"lij_IT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"li_NL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"lo_LA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"lt_LT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"lv_LV.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mai_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mg_MG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mhr_RU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mi_NZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mk_MK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ml_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mn_MN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mr_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ms_MY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"mt_MT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"my_MM.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nb_NO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nds_DE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nds_NL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ne_NP.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nl_AW.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nl_BE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nl_NL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nn_NO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nr_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"nso_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"oc_FR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"om_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"om_KE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"or_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"os_RU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"pa_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"pap_AN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"pa_PK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"pl_PL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ps_AF.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"pt_BR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"pt_PT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ro_RO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ru_RU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ru_UA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"rw_RW.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sa_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sc_IT.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sd_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"se_NO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"shs_CA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sid_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"si_LK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sk_SK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sl_SI.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"so_DJ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"so_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"so_KE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"so_SO.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sq_AL.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sq_MK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sr_ME.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sr_RS.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ss_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"st_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sv_FI.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sv_SE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sw_KE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"sw_TZ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ta_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ta_LK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"te_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tg_TJ.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"th_TH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ti_ER.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ti_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tig_ER.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tk_TM.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tl_PH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tn_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tr_CY.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tr_TR.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ts_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"tt_RU.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ug_CN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"uk_UA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"unm_US.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ur_IN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ur_PK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"ve_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"vi_VN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"wa_BE.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"wae_CH.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"wal_ET.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"wo_SN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"xh_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"yi_US.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"yo_NG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"yue_HK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"zh_CN.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"zh_HK.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"zh_SG.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"zh_TW.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
		"zu_ZA.UTF-8": {
			"SQL_ASCII",
			"UTF8",
		},
	}
}
