// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

type erContext struct {
	rce *RunControlEntry
	// need GRS
	// need current memory (but that is obtainable through GRS, I think)
}

type erHandler interface {
	Invoke(ctx *erContext)
}

var erHandlers = map[uint]erHandler{
	01:    nil, // IO$
	02:    nil, // IOI$
	03:    nil, // IOW$
	04:    nil, // EDJS$
	06:    nil, // WAIT$
	07:    nil, // WANY$
	010:   nil, // COM$
	011:   &erEXITHandler{},
	012:   &erABORTHandler{},
	013:   nil, // FORK$
	014:   nil, // TFORK$
	015:   nil, // READ$
	016:   nil, // PRINT$
	017:   &erCSFHandler{},
	022:   nil, // DATE$
	023:   nil, // TIME$
	024:   nil, // IOWI$
	025:   nil, // IOXI$
	026:   &erEABTHandler{},
	027:   nil, // II$
	030:   nil, // ABSAD$
	032:   nil, // FITEM$
	033:   nil, // INT$
	034:   nil, // IDENT$
	035:   nil, // CRTN$
	037:   nil, // WALL$
	040:   &erERRHandler{},
	041:   nil, // MCT$
	042:   nil, // READA$
	043:   nil, // MCORE$
	044:   nil, // LCORE$
	054:   nil, // TDATE$
	060:   nil, // TWAIT$
	061:   nil, // RT$
	062:   nil, // NRT$
	063:   nil, // OPT$
	064:   nil, // PCT$
	065:   nil, // SETC$
	066:   nil, // COND$
	067:   nil, // UNLCK$
	070:   nil, // APRINT$
	071:   nil, // APRNTA$
	072:   nil, // APUNCH$
	073:   nil, // APNCHA$
	074:   nil, // APRTCN$
	075:   nil, // APCHCN$
	076:   nil, // APRTCA$
	077:   nil, // APCHCA$
	0100:  nil, // CEND$
	0101:  nil, // IALL$
	0102:  nil, // TREAD$
	0103:  nil, // SWAIT$
	0104:  nil, // PFI$
	0105:  nil, // PFS$
	0106:  nil, // PFD$
	0107:  nil, // PFUWL$
	0110:  nil, // PFWL$
	0111:  nil, // LOAD$
	0112:  nil, // RSI$
	0113:  nil, // TSQCL$
	0114:  nil, // FACIL$
	0115:  nil, // BDSPT$
	0116:  nil, // INFO$
	0117:  nil, // CQUE$
	0120:  nil, // TRMRG$
	0121:  nil, // TSQRG$
	0122:  nil, // CTSQ$
	0123:  nil, // CTS$
	0124:  nil, // CTSA$
	0125:  nil, // MSCON$
	0126:  nil, // SNAP$
	0130:  nil, // PUNCH$
	0134:  nil, // AWAIT$
	0135:  nil, // TSWAP$
	0136:  nil, // TINTL$
	0137:  nil, // PRTCN$
	0140:  &erACSFHandler{},
	0141:  nil, // TOUT$
	0142:  nil, // TLBL$
	0143:  nil, // FACIT$
	0144:  nil, // PRNTA$
	0145:  nil, // PNCHA$
	0146:  nil, // NAME$
	0147:  nil, // ACT$
	0150:  nil, // DACT$
	0153:  nil, // CLIST$
	0155:  nil, // PRTCA$
	0156:  nil, // SETBP$
	0157:  nil, // PSR$
	0160:  nil, // BANK$
	0161:  nil, // ADED$
	0163:  nil, // ACCNT$
	0164:  nil, // PCHCN$
	0165:  nil, // PCHCA$
	0166:  nil, // AREAD$
	0167:  nil, // AREADA$
	0170:  nil, // ATREAD$
	0176:  nil, // SYSBAL$
	0200:  nil, // SYMB$
	0202:  nil, // ERRPR$
	0207:  nil, // LEVEL$
	0210:  nil, // LOG$
	0212:  nil, // CREG$
	0213:  nil, // SREG$
	0214:  nil, // SUVAL$
	0215:  nil, // SUMOD$
	0216:  nil, // STAB$
	0222:  nil, // SDEL$
	0223:  nil, // SPRNT$
	0225:  nil, // SABORT$
	0233:  nil, // DMGC$
	0234:  nil, // ERCVS$
	0235:  nil, // MQF$
	0236:  nil, // SC$QR
	0237:  nil, // DMABT$
	0241:  nil, // AUDIT$
	0242:  nil, // SYMINFO$
	0243:  nil, // SMOQUE$
	0244:  nil, // KEYIN$
	0246:  nil, // HMDBIT$
	0247:  &erCSIHandler{},
	0250:  nil, // CONFIG$
	0251:  nil, // TRTIM$
	0252:  nil, // ERTRAP$
	0253:  nil, // REGRTN$
	0254:  nil, // REGREP$
	0255:  nil, // TRAPRTN$
	0263:  nil, // TRON$
	0264:  nil, // DWTIME$
	0270:  nil, // AP$KEY
	0271:  nil, // AT$KEY
	0272:  nil, // SYSLOG$
	0273:  nil, // MODPS$
	0274:  nil, // TERMRUN$
	0277:  nil, // QECL$
	0300:  nil, // DQECL$
	0303:  nil, // SATTCP$
	0304:  nil, // SCDTL$
	0305:  nil, // SCDTA$
	0307:  nil, // TVSLBL$
	0312:  nil, // SCLDT$
	0313:  nil, // SCOMCNV$
	0314:  nil, // H2CON$
	0320:  nil, // SYS$TIME
	0321:  nil, // MODSWTIME$
	0322:  nil, // TIMECONFIG$
	0323:  nil, // TIMEBYINDEX$
	02004: nil, // RT$INT
	02005: nil, // RT$OUT
	02006: nil, // CMS$REG
	02011: nil, // CA$ASG
	02012: nil, // CA$REL
	02021: nil, // CR$ELG
	02030: nil, // AC$NIT
	02031: nil, // VT$RD
	02041: nil, // VT$CHG
	02042: nil, // VT$PUR
	02044: nil, // TP$APL
	02046: nil, // TF$KEY
	02050: nil, // DM$FAC
	02051: nil, // DM$IO
	02052: nil, // DM$IOW
	02053: nil, // DM$WT
	02056: nil, // FLAGBOX
	02060: nil, // RT$PSI
	02061: nil, // RT$PSD
	02064: nil, // TPLIB$
	02065: nil, // XFR$
	02066: nil, // CALL$
	02067: nil, // RTN$
	02070: nil, // TCORE$
	02071: nil, // XRS$
	02074: nil, // CO$MIT
	02075: nil, // RL$BAK
	02101: nil, // RT$PSS
	02102: nil, // RT$PID
	02103: nil, // SEXEM$
	02104: nil, // TIP$Q
	02106: nil, // QI$NIT
	02107: nil, // QI$CON
	02110: nil, // QI$DIS
	02111: nil, // TIP$TA
	02112: nil, // TIP$TC
	02113: nil, // TIP$ID
	02114: nil, // MCABT$
	02115: nil, // MSGN$
	02117: nil, // PERF$
	02120: nil, // TIP$XMIT
	02130: nil, // TIP$SM
	02131: nil, // TIP$TALK
	02132: nil, // SC$SR
	02133: nil, // TM$SET
}

type erABORTHandler struct{}
type erACSFHandler struct{}
type erCSFHandler struct{}
type erCSIHandler struct{}
type erEABTHandler struct{}
type erERRHandler struct{}
type erEXITHandler struct{}

func (h *erABORTHandler) Invoke(ctx *erContext) {

}

func (h *erACSFHandler) Invoke(ctx *erContext) {
	//L A0,(image-length,image-address)
}

func (h *erCSFHandler) Invoke(ctx *erContext) {
	//L A0,(image-length,image-address)
}

func (h *erCSIHandler) Invoke(ctx *erContext) {
	//L A0,(image-length,image-address)
}

func (h *erEABTHandler) Invoke(ctx *erContext) {

}

func (h *erERRHandler) Invoke(ctx *erContext) {

}

func (h *erEXITHandler) Invoke(ctx *erContext) {

}
