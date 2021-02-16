// Code generated by go run make_schema.go. DO NOT EDIT.

package matroska

import (
	"time"

	"github.com/coding-socks/ebml"
)

var DocType = ebml.NewDefinition("18538067", ebml.TypeMaster, "Segment", nil, []ebml.Definition{
	ebml.NewDefinition("114D9B74", ebml.TypeMaster, "SeekHead", nil, []ebml.Definition{
		ebml.NewDefinition("4DBB", ebml.TypeMaster, "Seek", nil, []ebml.Definition{
			ebml.NewDefinition("53AB", ebml.TypeBinary, "SeekID", nil, nil),
			ebml.NewDefinition("53AC", ebml.TypeUinteger, "SeekPosition", nil, nil),
		}),
	}),
	ebml.NewDefinition("1549A966", ebml.TypeMaster, "Info", nil, []ebml.Definition{
		ebml.NewDefinition("73A4", ebml.TypeBinary, "SegmentUID", nil, nil),
		ebml.NewDefinition("7384", ebml.TypeUTF8, "SegmentFilename", nil, nil),
		ebml.NewDefinition("3CB923", ebml.TypeBinary, "PrevUID", nil, nil),
		ebml.NewDefinition("3C83AB", ebml.TypeUTF8, "PrevFilename", nil, nil),
		ebml.NewDefinition("3EB923", ebml.TypeBinary, "NextUID", nil, nil),
		ebml.NewDefinition("3E83BB", ebml.TypeUTF8, "NextFilename", nil, nil),
		ebml.NewDefinition("4444", ebml.TypeBinary, "SegmentFamily", nil, nil),
		ebml.NewDefinition("6924", ebml.TypeMaster, "ChapterTranslate", nil, []ebml.Definition{
			ebml.NewDefinition("69FC", ebml.TypeUinteger, "ChapterTranslateEditionUID", nil, nil),
			ebml.NewDefinition("69BF", ebml.TypeUinteger, "ChapterTranslateCodec", nil, nil),
			ebml.NewDefinition("69A5", ebml.TypeBinary, "ChapterTranslateID", nil, nil),
		}),
		ebml.NewDefinition("2AD7B1", ebml.TypeUinteger, "TimestampScale", 1000000, nil),
		ebml.NewDefinition("4489", ebml.TypeFloat, "Duration", nil, nil),
		ebml.NewDefinition("4461", ebml.TypeDate, "DateUTC", nil, nil),
		ebml.NewDefinition("7BA9", ebml.TypeUTF8, "Title", nil, nil),
		ebml.NewDefinition("4D80", ebml.TypeUTF8, "MuxingApp", nil, nil),
		ebml.NewDefinition("5741", ebml.TypeUTF8, "WritingApp", nil, nil),
	}),
	ebml.NewDefinition("1F43B675", ebml.TypeMaster, "Cluster", nil, []ebml.Definition{
		ebml.NewDefinition("E7", ebml.TypeUinteger, "Timestamp", nil, nil),
		ebml.NewDefinition("5854", ebml.TypeMaster, "SilentTracks", nil, []ebml.Definition{
			ebml.NewDefinition("58D7", ebml.TypeUinteger, "SilentTrackNumber", nil, nil),
		}),
		ebml.NewDefinition("A7", ebml.TypeUinteger, "Position", nil, nil),
		ebml.NewDefinition("AB", ebml.TypeUinteger, "PrevSize", nil, nil),
		ebml.NewDefinition("A3", ebml.TypeBinary, "SimpleBlock", nil, nil),
		ebml.NewDefinition("A0", ebml.TypeMaster, "BlockGroup", nil, []ebml.Definition{
			ebml.NewDefinition("A1", ebml.TypeBinary, "Block", nil, nil),
			ebml.NewDefinition("A2", ebml.TypeBinary, "BlockVirtual", nil, nil),
			ebml.NewDefinition("75A1", ebml.TypeMaster, "BlockAdditions", nil, []ebml.Definition{
				ebml.NewDefinition("A6", ebml.TypeMaster, "BlockMore", nil, []ebml.Definition{
					ebml.NewDefinition("EE", ebml.TypeUinteger, "BlockAddID", 1, nil),
					ebml.NewDefinition("A5", ebml.TypeBinary, "BlockAdditional", nil, nil),
				}),
			}),
			ebml.NewDefinition("9B", ebml.TypeUinteger, "BlockDuration", nil, nil),
			ebml.NewDefinition("FA", ebml.TypeUinteger, "ReferencePriority", 0, nil),
			ebml.NewDefinition("FB", ebml.TypeInteger, "ReferenceBlock", nil, nil),
			ebml.NewDefinition("FD", ebml.TypeInteger, "ReferenceVirtual", nil, nil),
			ebml.NewDefinition("A4", ebml.TypeBinary, "CodecState", nil, nil),
			ebml.NewDefinition("75A2", ebml.TypeInteger, "DiscardPadding", nil, nil),
			ebml.NewDefinition("8E", ebml.TypeMaster, "Slices", nil, []ebml.Definition{
				ebml.NewDefinition("E8", ebml.TypeMaster, "TimeSlice", nil, []ebml.Definition{
					ebml.NewDefinition("CC", ebml.TypeUinteger, "LaceNumber", 0, nil),
					ebml.NewDefinition("CD", ebml.TypeUinteger, "FrameNumber", 0, nil),
					ebml.NewDefinition("CB", ebml.TypeUinteger, "BlockAdditionID", 0, nil),
					ebml.NewDefinition("CE", ebml.TypeUinteger, "Delay", 0, nil),
					ebml.NewDefinition("CF", ebml.TypeUinteger, "SliceDuration", 0, nil),
				}),
			}),
			ebml.NewDefinition("C8", ebml.TypeMaster, "ReferenceFrame", nil, []ebml.Definition{
				ebml.NewDefinition("C9", ebml.TypeUinteger, "ReferenceOffset", nil, nil),
				ebml.NewDefinition("CA", ebml.TypeUinteger, "ReferenceTimestamp", nil, nil),
			}),
		}),
		ebml.NewDefinition("AF", ebml.TypeBinary, "EncryptedBlock", nil, nil),
	}),
	ebml.NewDefinition("1654AE6B", ebml.TypeMaster, "Tracks", nil, []ebml.Definition{
		ebml.NewDefinition("AE", ebml.TypeMaster, "TrackEntry", nil, []ebml.Definition{
			ebml.NewDefinition("D7", ebml.TypeUinteger, "TrackNumber", nil, nil),
			ebml.NewDefinition("73C5", ebml.TypeUinteger, "TrackUID", nil, nil),
			ebml.NewDefinition("83", ebml.TypeUinteger, "TrackType", nil, nil),
			ebml.NewDefinition("B9", ebml.TypeUinteger, "FlagEnabled", 1, nil),
			ebml.NewDefinition("88", ebml.TypeUinteger, "FlagDefault", 1, nil),
			ebml.NewDefinition("55AA", ebml.TypeUinteger, "FlagForced", 0, nil),
			ebml.NewDefinition("9C", ebml.TypeUinteger, "FlagLacing", 1, nil),
			ebml.NewDefinition("6DE7", ebml.TypeUinteger, "MinCache", 0, nil),
			ebml.NewDefinition("6DF8", ebml.TypeUinteger, "MaxCache", nil, nil),
			ebml.NewDefinition("23E383", ebml.TypeUinteger, "DefaultDuration", nil, nil),
			ebml.NewDefinition("234E7A", ebml.TypeUinteger, "DefaultDecodedFieldDuration", nil, nil),
			ebml.NewDefinition("23314F", ebml.TypeFloat, "TrackTimestampScale", 0x1p+0, nil),
			ebml.NewDefinition("537F", ebml.TypeInteger, "TrackOffset", 0, nil),
			ebml.NewDefinition("55EE", ebml.TypeUinteger, "MaxBlockAdditionID", 0, nil),
			ebml.NewDefinition("41E4", ebml.TypeMaster, "BlockAdditionMapping", nil, []ebml.Definition{
				ebml.NewDefinition("41F0", ebml.TypeUinteger, "BlockAddIDValue", nil, nil),
				ebml.NewDefinition("41A4", ebml.TypeString, "BlockAddIDName", nil, nil),
				ebml.NewDefinition("41E7", ebml.TypeUinteger, "BlockAddIDType", 0, nil),
				ebml.NewDefinition("41ED", ebml.TypeBinary, "BlockAddIDExtraData", nil, nil),
			}),
			ebml.NewDefinition("536E", ebml.TypeUTF8, "Name", nil, nil),
			ebml.NewDefinition("22B59C", ebml.TypeString, "Language", "eng", nil),
			ebml.NewDefinition("22B59D", ebml.TypeString, "LanguageIETF", nil, nil),
			ebml.NewDefinition("86", ebml.TypeString, "CodecID", nil, nil),
			ebml.NewDefinition("63A2", ebml.TypeBinary, "CodecPrivate", nil, nil),
			ebml.NewDefinition("258688", ebml.TypeUTF8, "CodecName", nil, nil),
			ebml.NewDefinition("7446", ebml.TypeUinteger, "AttachmentLink", nil, nil),
			ebml.NewDefinition("3A9697", ebml.TypeUTF8, "CodecSettings", nil, nil),
			ebml.NewDefinition("3B4040", ebml.TypeString, "CodecInfoURL", nil, nil),
			ebml.NewDefinition("26B240", ebml.TypeString, "CodecDownloadURL", nil, nil),
			ebml.NewDefinition("AA", ebml.TypeUinteger, "CodecDecodeAll", 1, nil),
			ebml.NewDefinition("6FAB", ebml.TypeUinteger, "TrackOverlay", nil, nil),
			ebml.NewDefinition("56AA", ebml.TypeUinteger, "CodecDelay", 0, nil),
			ebml.NewDefinition("56BB", ebml.TypeUinteger, "SeekPreRoll", 0, nil),
			ebml.NewDefinition("6624", ebml.TypeMaster, "TrackTranslate", nil, []ebml.Definition{
				ebml.NewDefinition("66FC", ebml.TypeUinteger, "TrackTranslateEditionUID", nil, nil),
				ebml.NewDefinition("66BF", ebml.TypeUinteger, "TrackTranslateCodec", nil, nil),
				ebml.NewDefinition("66A5", ebml.TypeBinary, "TrackTranslateTrackID", nil, nil),
			}),
			ebml.NewDefinition("E0", ebml.TypeMaster, "Video", nil, []ebml.Definition{
				ebml.NewDefinition("9A", ebml.TypeUinteger, "FlagInterlaced", 0, nil),
				ebml.NewDefinition("9D", ebml.TypeUinteger, "FieldOrder", 2, nil),
				ebml.NewDefinition("53B8", ebml.TypeUinteger, "StereoMode", 0, nil),
				ebml.NewDefinition("53C0", ebml.TypeUinteger, "AlphaMode", 0, nil),
				ebml.NewDefinition("53B9", ebml.TypeUinteger, "OldStereoMode", nil, nil),
				ebml.NewDefinition("B0", ebml.TypeUinteger, "PixelWidth", nil, nil),
				ebml.NewDefinition("BA", ebml.TypeUinteger, "PixelHeight", nil, nil),
				ebml.NewDefinition("54AA", ebml.TypeUinteger, "PixelCropBottom", 0, nil),
				ebml.NewDefinition("54BB", ebml.TypeUinteger, "PixelCropTop", 0, nil),
				ebml.NewDefinition("54CC", ebml.TypeUinteger, "PixelCropLeft", 0, nil),
				ebml.NewDefinition("54DD", ebml.TypeUinteger, "PixelCropRight", 0, nil),
				ebml.NewDefinition("54B0", ebml.TypeUinteger, "DisplayWidth", nil, nil),
				ebml.NewDefinition("54BA", ebml.TypeUinteger, "DisplayHeight", nil, nil),
				ebml.NewDefinition("54B2", ebml.TypeUinteger, "DisplayUnit", 0, nil),
				ebml.NewDefinition("54B3", ebml.TypeUinteger, "AspectRatioType", 0, nil),
				ebml.NewDefinition("2EB524", ebml.TypeBinary, "ColourSpace", nil, nil),
				ebml.NewDefinition("2FB523", ebml.TypeFloat, "GammaValue", nil, nil),
				ebml.NewDefinition("2383E3", ebml.TypeFloat, "FrameRate", nil, nil),
				ebml.NewDefinition("55B0", ebml.TypeMaster, "Colour", nil, []ebml.Definition{
					ebml.NewDefinition("55B1", ebml.TypeUinteger, "MatrixCoefficients", 2, nil),
					ebml.NewDefinition("55B2", ebml.TypeUinteger, "BitsPerChannel", 0, nil),
					ebml.NewDefinition("55B3", ebml.TypeUinteger, "ChromaSubsamplingHorz", nil, nil),
					ebml.NewDefinition("55B4", ebml.TypeUinteger, "ChromaSubsamplingVert", nil, nil),
					ebml.NewDefinition("55B5", ebml.TypeUinteger, "CbSubsamplingHorz", nil, nil),
					ebml.NewDefinition("55B6", ebml.TypeUinteger, "CbSubsamplingVert", nil, nil),
					ebml.NewDefinition("55B7", ebml.TypeUinteger, "ChromaSitingHorz", 0, nil),
					ebml.NewDefinition("55B8", ebml.TypeUinteger, "ChromaSitingVert", 0, nil),
					ebml.NewDefinition("55B9", ebml.TypeUinteger, "Range", 0, nil),
					ebml.NewDefinition("55BA", ebml.TypeUinteger, "TransferCharacteristics", 2, nil),
					ebml.NewDefinition("55BB", ebml.TypeUinteger, "Primaries", 2, nil),
					ebml.NewDefinition("55BC", ebml.TypeUinteger, "MaxCLL", nil, nil),
					ebml.NewDefinition("55BD", ebml.TypeUinteger, "MaxFALL", nil, nil),
					ebml.NewDefinition("55D0", ebml.TypeMaster, "MasteringMetadata", nil, []ebml.Definition{
						ebml.NewDefinition("55D1", ebml.TypeFloat, "PrimaryRChromaticityX", nil, nil),
						ebml.NewDefinition("55D2", ebml.TypeFloat, "PrimaryRChromaticityY", nil, nil),
						ebml.NewDefinition("55D3", ebml.TypeFloat, "PrimaryGChromaticityX", nil, nil),
						ebml.NewDefinition("55D4", ebml.TypeFloat, "PrimaryGChromaticityY", nil, nil),
						ebml.NewDefinition("55D5", ebml.TypeFloat, "PrimaryBChromaticityX", nil, nil),
						ebml.NewDefinition("55D6", ebml.TypeFloat, "PrimaryBChromaticityY", nil, nil),
						ebml.NewDefinition("55D7", ebml.TypeFloat, "WhitePointChromaticityX", nil, nil),
						ebml.NewDefinition("55D8", ebml.TypeFloat, "WhitePointChromaticityY", nil, nil),
						ebml.NewDefinition("55D9", ebml.TypeFloat, "LuminanceMax", nil, nil),
						ebml.NewDefinition("55DA", ebml.TypeFloat, "LuminanceMin", nil, nil),
					}),
				}),
				ebml.NewDefinition("7670", ebml.TypeMaster, "Projection", nil, []ebml.Definition{
					ebml.NewDefinition("7671", ebml.TypeUinteger, "ProjectionType", 0, nil),
					ebml.NewDefinition("7672", ebml.TypeBinary, "ProjectionPrivate", nil, nil),
					ebml.NewDefinition("7673", ebml.TypeFloat, "ProjectionPoseYaw", 0x0p+0, nil),
					ebml.NewDefinition("7674", ebml.TypeFloat, "ProjectionPosePitch", 0x0p+0, nil),
					ebml.NewDefinition("7675", ebml.TypeFloat, "ProjectionPoseRoll", 0x0p+0, nil),
				}),
			}),
			ebml.NewDefinition("E1", ebml.TypeMaster, "Audio", nil, []ebml.Definition{
				ebml.NewDefinition("B5", ebml.TypeFloat, "SamplingFrequency", 0x1.f4p+12, nil),
				ebml.NewDefinition("78B5", ebml.TypeFloat, "OutputSamplingFrequency", nil, nil),
				ebml.NewDefinition("9F", ebml.TypeUinteger, "Channels", 1, nil),
				ebml.NewDefinition("7D7B", ebml.TypeBinary, "ChannelPositions", nil, nil),
				ebml.NewDefinition("6264", ebml.TypeUinteger, "BitDepth", nil, nil),
			}),
			ebml.NewDefinition("E2", ebml.TypeMaster, "TrackOperation", nil, []ebml.Definition{
				ebml.NewDefinition("E3", ebml.TypeMaster, "TrackCombinePlanes", nil, []ebml.Definition{
					ebml.NewDefinition("E4", ebml.TypeMaster, "TrackPlane", nil, []ebml.Definition{
						ebml.NewDefinition("E5", ebml.TypeUinteger, "TrackPlaneUID", nil, nil),
						ebml.NewDefinition("E6", ebml.TypeUinteger, "TrackPlaneType", nil, nil),
					}),
				}),
				ebml.NewDefinition("E9", ebml.TypeMaster, "TrackJoinBlocks", nil, []ebml.Definition{
					ebml.NewDefinition("ED", ebml.TypeUinteger, "TrackJoinUID", nil, nil),
				}),
			}),
			ebml.NewDefinition("C0", ebml.TypeUinteger, "TrickTrackUID", nil, nil),
			ebml.NewDefinition("C1", ebml.TypeBinary, "TrickTrackSegmentUID", nil, nil),
			ebml.NewDefinition("C6", ebml.TypeUinteger, "TrickTrackFlag", 0, nil),
			ebml.NewDefinition("C7", ebml.TypeUinteger, "TrickMasterTrackUID", nil, nil),
			ebml.NewDefinition("C4", ebml.TypeBinary, "TrickMasterTrackSegmentUID", nil, nil),
			ebml.NewDefinition("6D80", ebml.TypeMaster, "ContentEncodings", nil, []ebml.Definition{
				ebml.NewDefinition("6240", ebml.TypeMaster, "ContentEncoding", nil, []ebml.Definition{
					ebml.NewDefinition("5031", ebml.TypeUinteger, "ContentEncodingOrder", 0, nil),
					ebml.NewDefinition("5032", ebml.TypeUinteger, "ContentEncodingScope", 1, nil),
					ebml.NewDefinition("5033", ebml.TypeUinteger, "ContentEncodingType", 0, nil),
					ebml.NewDefinition("5034", ebml.TypeMaster, "ContentCompression", nil, []ebml.Definition{
						ebml.NewDefinition("4254", ebml.TypeUinteger, "ContentCompAlgo", 0, nil),
						ebml.NewDefinition("4255", ebml.TypeBinary, "ContentCompSettings", nil, nil),
					}),
					ebml.NewDefinition("5035", ebml.TypeMaster, "ContentEncryption", nil, []ebml.Definition{
						ebml.NewDefinition("47E1", ebml.TypeUinteger, "ContentEncAlgo", 0, nil),
						ebml.NewDefinition("47E2", ebml.TypeBinary, "ContentEncKeyID", nil, nil),
						ebml.NewDefinition("47E7", ebml.TypeMaster, "ContentEncAESSettings", nil, []ebml.Definition{
							ebml.NewDefinition("47E8", ebml.TypeUinteger, "AESSettingsCipherMode", nil, nil),
						}),
						ebml.NewDefinition("47E3", ebml.TypeBinary, "ContentSignature", nil, nil),
						ebml.NewDefinition("47E4", ebml.TypeBinary, "ContentSigKeyID", nil, nil),
						ebml.NewDefinition("47E5", ebml.TypeUinteger, "ContentSigAlgo", 0, nil),
						ebml.NewDefinition("47E6", ebml.TypeUinteger, "ContentSigHashAlgo", 0, nil),
					}),
				}),
			}),
		}),
	}),
	ebml.NewDefinition("1C53BB6B", ebml.TypeMaster, "Cues", nil, []ebml.Definition{
		ebml.NewDefinition("BB", ebml.TypeMaster, "CuePoint", nil, []ebml.Definition{
			ebml.NewDefinition("B3", ebml.TypeUinteger, "CueTime", nil, nil),
			ebml.NewDefinition("B7", ebml.TypeMaster, "CueTrackPositions", nil, []ebml.Definition{
				ebml.NewDefinition("F7", ebml.TypeUinteger, "CueTrack", nil, nil),
				ebml.NewDefinition("F1", ebml.TypeUinteger, "CueClusterPosition", nil, nil),
				ebml.NewDefinition("F0", ebml.TypeUinteger, "CueRelativePosition", nil, nil),
				ebml.NewDefinition("B2", ebml.TypeUinteger, "CueDuration", nil, nil),
				ebml.NewDefinition("5378", ebml.TypeUinteger, "CueBlockNumber", 1, nil),
				ebml.NewDefinition("EA", ebml.TypeUinteger, "CueCodecState", 0, nil),
				ebml.NewDefinition("DB", ebml.TypeMaster, "CueReference", nil, []ebml.Definition{
					ebml.NewDefinition("96", ebml.TypeUinteger, "CueRefTime", nil, nil),
					ebml.NewDefinition("97", ebml.TypeUinteger, "CueRefCluster", nil, nil),
					ebml.NewDefinition("535F", ebml.TypeUinteger, "CueRefNumber", 1, nil),
					ebml.NewDefinition("EB", ebml.TypeUinteger, "CueRefCodecState", 0, nil),
				}),
			}),
		}),
	}),
	ebml.NewDefinition("1941A469", ebml.TypeMaster, "Attachments", nil, []ebml.Definition{
		ebml.NewDefinition("61A7", ebml.TypeMaster, "AttachedFile", nil, []ebml.Definition{
			ebml.NewDefinition("467E", ebml.TypeUTF8, "FileDescription", nil, nil),
			ebml.NewDefinition("466E", ebml.TypeUTF8, "FileName", nil, nil),
			ebml.NewDefinition("4660", ebml.TypeString, "FileMimeType", nil, nil),
			ebml.NewDefinition("465C", ebml.TypeBinary, "FileData", nil, nil),
			ebml.NewDefinition("46AE", ebml.TypeUinteger, "FileUID", nil, nil),
			ebml.NewDefinition("4675", ebml.TypeBinary, "FileReferral", nil, nil),
			ebml.NewDefinition("4661", ebml.TypeUinteger, "FileUsedStartTime", nil, nil),
			ebml.NewDefinition("4662", ebml.TypeUinteger, "FileUsedEndTime", nil, nil),
		}),
	}),
	ebml.NewDefinition("1043A770", ebml.TypeMaster, "Chapters", nil, []ebml.Definition{
		ebml.NewDefinition("45B9", ebml.TypeMaster, "EditionEntry", nil, []ebml.Definition{
			ebml.NewDefinition("45BC", ebml.TypeUinteger, "EditionUID", nil, nil),
			ebml.NewDefinition("45BD", ebml.TypeUinteger, "EditionFlagHidden", 0, nil),
			ebml.NewDefinition("45DB", ebml.TypeUinteger, "EditionFlagDefault", 0, nil),
			ebml.NewDefinition("45DD", ebml.TypeUinteger, "EditionFlagOrdered", 0, nil),
			ebml.NewDefinition("B6", ebml.TypeMaster, "ChapterAtom", nil, []ebml.Definition{
				ebml.NewDefinition("73C4", ebml.TypeUinteger, "ChapterUID", nil, nil),
				ebml.NewDefinition("5654", ebml.TypeUTF8, "ChapterStringUID", nil, nil),
				ebml.NewDefinition("91", ebml.TypeUinteger, "ChapterTimeStart", nil, nil),
				ebml.NewDefinition("92", ebml.TypeUinteger, "ChapterTimeEnd", nil, nil),
				ebml.NewDefinition("98", ebml.TypeUinteger, "ChapterFlagHidden", 0, nil),
				ebml.NewDefinition("4598", ebml.TypeUinteger, "ChapterFlagEnabled", 1, nil),
				ebml.NewDefinition("6E67", ebml.TypeBinary, "ChapterSegmentUID", nil, nil),
				ebml.NewDefinition("6EBC", ebml.TypeUinteger, "ChapterSegmentEditionUID", nil, nil),
				ebml.NewDefinition("63C3", ebml.TypeUinteger, "ChapterPhysicalEquiv", nil, nil),
				ebml.NewDefinition("8F", ebml.TypeMaster, "ChapterTrack", nil, []ebml.Definition{
					ebml.NewDefinition("89", ebml.TypeUinteger, "ChapterTrackUID", nil, nil),
				}),
				ebml.NewDefinition("80", ebml.TypeMaster, "ChapterDisplay", nil, []ebml.Definition{
					ebml.NewDefinition("85", ebml.TypeUTF8, "ChapString", nil, nil),
					ebml.NewDefinition("437C", ebml.TypeString, "ChapLanguage", "eng", nil),
					ebml.NewDefinition("437D", ebml.TypeString, "ChapLanguageIETF", nil, nil),
					ebml.NewDefinition("437E", ebml.TypeString, "ChapCountry", nil, nil),
				}),
				ebml.NewDefinition("6944", ebml.TypeMaster, "ChapProcess", nil, []ebml.Definition{
					ebml.NewDefinition("6955", ebml.TypeUinteger, "ChapProcessCodecID", 0, nil),
					ebml.NewDefinition("450D", ebml.TypeBinary, "ChapProcessPrivate", nil, nil),
					ebml.NewDefinition("6911", ebml.TypeMaster, "ChapProcessCommand", nil, []ebml.Definition{
						ebml.NewDefinition("6922", ebml.TypeUinteger, "ChapProcessTime", nil, nil),
						ebml.NewDefinition("6933", ebml.TypeBinary, "ChapProcessData", nil, nil),
					}),
				}),
			}),
		}),
	}),
	ebml.NewDefinition("1254C367", ebml.TypeMaster, "Tags", nil, []ebml.Definition{
		ebml.NewDefinition("7373", ebml.TypeMaster, "Tag", nil, []ebml.Definition{
			ebml.NewDefinition("63C0", ebml.TypeMaster, "Targets", nil, []ebml.Definition{
				ebml.NewDefinition("68CA", ebml.TypeUinteger, "TargetTypeValue", 50, nil),
				ebml.NewDefinition("63CA", ebml.TypeString, "TargetType", nil, nil),
				ebml.NewDefinition("63C5", ebml.TypeUinteger, "TagTrackUID", 0, nil),
				ebml.NewDefinition("63C9", ebml.TypeUinteger, "TagEditionUID", 0, nil),
				ebml.NewDefinition("63C4", ebml.TypeUinteger, "TagChapterUID", 0, nil),
				ebml.NewDefinition("63C6", ebml.TypeUinteger, "TagAttachmentUID", 0, nil),
			}),
			ebml.NewDefinition("67C8", ebml.TypeMaster, "SimpleTag", nil, []ebml.Definition{
				ebml.NewDefinition("45A3", ebml.TypeUTF8, "TagName", nil, nil),
				ebml.NewDefinition("447A", ebml.TypeString, "TagLanguage", "und", nil),
				ebml.NewDefinition("447B", ebml.TypeString, "TagLanguageIETF", nil, nil),
				ebml.NewDefinition("4484", ebml.TypeUinteger, "TagDefault", 1, nil),
				ebml.NewDefinition("4487", ebml.TypeUTF8, "TagString", nil, nil),
				ebml.NewDefinition("4485", ebml.TypeBinary, "TagBinary", nil, nil),
			}),
		}),
	}),
})

type Document struct {
	EBML    EBML
	Segment Segment
}

type EBML struct {
	EBMLVersion        uint
	EBMLReadVersion    uint
	EBMLMaxIDLength    uint
	EBMLMaxSizeLength  uint
	DocType            string
	DocTypeVersion     uint
	DocTypeReadVersion uint
	DocTypeExtension   []DocTypeExtension
}

type DocTypeExtension struct {
	DocTypeExtensionName    string
	DocTypeExtensionVersion uint
}

type Segment struct {
	SeekHead    []SeekHead
	Info        Info
	Cluster     []Cluster
	Tracks      Tracks
	Cues        Cues
	Attachments Attachments
	Chapters    Chapters
	Tags        []Tags
}

type SeekHead struct {
	Seek []Seek
}

type Seek struct {
	SeekID       []byte
	SeekPosition uint
}

type Info struct {
	SegmentUID       []byte
	SegmentFilename  string
	PrevUID          []byte
	PrevFilename     string
	NextUID          []byte
	NextFilename     string
	SegmentFamily    [][]byte
	ChapterTranslate []ChapterTranslate
	TimestampScale   uint
	Duration         float64
	DateUTC          time.Time
	Title            string
	MuxingApp        string
	WritingApp       string
}

type ChapterTranslate struct {
	ChapterTranslateEditionUID []uint
	ChapterTranslateCodec      uint
	ChapterTranslateID         []byte
}

type Cluster struct {
	Timestamp      uint
	SilentTracks   SilentTracks
	Position       uint
	PrevSize       uint
	SimpleBlock    [][]byte
	BlockGroup     []BlockGroup
	EncryptedBlock [][]byte
}

type SilentTracks struct {
	SilentTrackNumber []uint
}

type BlockGroup struct {
	Block             []byte
	BlockVirtual      []byte
	BlockAdditions    BlockAdditions
	BlockDuration     uint
	ReferencePriority uint
	ReferenceBlock    []int
	ReferenceVirtual  int
	CodecState        []byte
	DiscardPadding    int
	Slices            Slices
	ReferenceFrame    ReferenceFrame
}

type BlockAdditions struct {
	BlockMore []BlockMore
}

type BlockMore struct {
	BlockAddID      uint
	BlockAdditional []byte
}

type Slices struct {
	TimeSlice []TimeSlice
}

type TimeSlice struct {
	LaceNumber      uint
	FrameNumber     uint
	BlockAdditionID uint
	Delay           uint
	SliceDuration   uint
}

type ReferenceFrame struct {
	ReferenceOffset    uint
	ReferenceTimestamp uint
}

type Tracks struct {
	TrackEntry []TrackEntry
}

type TrackEntry struct {
	TrackNumber                 uint
	TrackUID                    uint
	TrackType                   uint
	FlagEnabled                 uint
	FlagDefault                 uint
	FlagForced                  uint
	FlagLacing                  uint
	MinCache                    uint
	MaxCache                    uint
	DefaultDuration             uint
	DefaultDecodedFieldDuration uint
	TrackTimestampScale         float64
	TrackOffset                 int
	MaxBlockAdditionID          uint
	BlockAdditionMapping        []BlockAdditionMapping
	Name                        string
	Language                    string
	LanguageIETF                string
	CodecID                     string
	CodecPrivate                []byte
	CodecName                   string
	AttachmentLink              uint
	CodecSettings               string
	CodecInfoURL                []string
	CodecDownloadURL            []string
	CodecDecodeAll              uint
	TrackOverlay                []uint
	CodecDelay                  uint
	SeekPreRoll                 uint
	TrackTranslate              []TrackTranslate
	Video                       Video
	Audio                       Audio
	TrackOperation              TrackOperation
	TrickTrackUID               uint
	TrickTrackSegmentUID        []byte
	TrickTrackFlag              uint
	TrickMasterTrackUID         uint
	TrickMasterTrackSegmentUID  []byte
	ContentEncodings            ContentEncodings
}

type BlockAdditionMapping struct {
	BlockAddIDValue     uint
	BlockAddIDName      string
	BlockAddIDType      uint
	BlockAddIDExtraData []byte
}

type TrackTranslate struct {
	TrackTranslateEditionUID []uint
	TrackTranslateCodec      uint
	TrackTranslateTrackID    []byte
}

type Video struct {
	FlagInterlaced  uint
	FieldOrder      uint
	StereoMode      uint
	AlphaMode       uint
	OldStereoMode   uint
	PixelWidth      uint
	PixelHeight     uint
	PixelCropBottom uint
	PixelCropTop    uint
	PixelCropLeft   uint
	PixelCropRight  uint
	DisplayWidth    uint
	DisplayHeight   uint
	DisplayUnit     uint
	AspectRatioType uint
	ColourSpace     []byte
	GammaValue      float64
	FrameRate       float64
	Colour          Colour
	Projection      Projection
}

type Colour struct {
	MatrixCoefficients      uint
	BitsPerChannel          uint
	ChromaSubsamplingHorz   uint
	ChromaSubsamplingVert   uint
	CbSubsamplingHorz       uint
	CbSubsamplingVert       uint
	ChromaSitingHorz        uint
	ChromaSitingVert        uint
	Range                   uint
	TransferCharacteristics uint
	Primaries               uint
	MaxCLL                  uint
	MaxFALL                 uint
	MasteringMetadata       MasteringMetadata
}

type MasteringMetadata struct {
	PrimaryRChromaticityX   float64
	PrimaryRChromaticityY   float64
	PrimaryGChromaticityX   float64
	PrimaryGChromaticityY   float64
	PrimaryBChromaticityX   float64
	PrimaryBChromaticityY   float64
	WhitePointChromaticityX float64
	WhitePointChromaticityY float64
	LuminanceMax            float64
	LuminanceMin            float64
}

type Projection struct {
	ProjectionType      uint
	ProjectionPrivate   []byte
	ProjectionPoseYaw   float64
	ProjectionPosePitch float64
	ProjectionPoseRoll  float64
}

type Audio struct {
	SamplingFrequency       float64
	OutputSamplingFrequency float64
	Channels                uint
	ChannelPositions        []byte
	BitDepth                uint
}

type TrackOperation struct {
	TrackCombinePlanes TrackCombinePlanes
	TrackJoinBlocks    TrackJoinBlocks
}

type TrackCombinePlanes struct {
	TrackPlane []TrackPlane
}

type TrackPlane struct {
	TrackPlaneUID  uint
	TrackPlaneType uint
}

type TrackJoinBlocks struct {
	TrackJoinUID []uint
}

type ContentEncodings struct {
	ContentEncoding []ContentEncoding
}

type ContentEncoding struct {
	ContentEncodingOrder uint
	ContentEncodingScope uint
	ContentEncodingType  uint
	ContentCompression   ContentCompression
	ContentEncryption    ContentEncryption
}

type ContentCompression struct {
	ContentCompAlgo     uint
	ContentCompSettings []byte
}

type ContentEncryption struct {
	ContentEncAlgo        uint
	ContentEncKeyID       []byte
	ContentEncAESSettings ContentEncAESSettings
	ContentSignature      []byte
	ContentSigKeyID       []byte
	ContentSigAlgo        uint
	ContentSigHashAlgo    uint
}

type ContentEncAESSettings struct {
	AESSettingsCipherMode uint
}

type Cues struct {
	CuePoint []CuePoint
}

type CuePoint struct {
	CueTime           uint
	CueTrackPositions []CueTrackPositions
}

type CueTrackPositions struct {
	CueTrack            uint
	CueClusterPosition  uint
	CueRelativePosition uint
	CueDuration         uint
	CueBlockNumber      uint
	CueCodecState       uint
	CueReference        []CueReference
}

type CueReference struct {
	CueRefTime       uint
	CueRefCluster    uint
	CueRefNumber     uint
	CueRefCodecState uint
}

type Attachments struct {
	AttachedFile []AttachedFile
}

type AttachedFile struct {
	FileDescription   string
	FileName          string
	FileMimeType      string
	FileData          []byte
	FileUID           uint
	FileReferral      []byte
	FileUsedStartTime uint
	FileUsedEndTime   uint
}

type Chapters struct {
	EditionEntry []EditionEntry
}

type EditionEntry struct {
	EditionUID         uint
	EditionFlagHidden  uint
	EditionFlagDefault uint
	EditionFlagOrdered uint
	ChapterAtom        []ChapterAtom
}

type ChapterAtom struct {
	ChapterUID               uint
	ChapterStringUID         string
	ChapterTimeStart         uint
	ChapterTimeEnd           uint
	ChapterFlagHidden        uint
	ChapterFlagEnabled       uint
	ChapterSegmentUID        []byte
	ChapterSegmentEditionUID uint
	ChapterPhysicalEquiv     uint
	ChapterTrack             ChapterTrack
	ChapterDisplay           []ChapterDisplay
	ChapProcess              []ChapProcess
}

type ChapterTrack struct {
	ChapterTrackUID []uint
}

type ChapterDisplay struct {
	ChapString       string
	ChapLanguage     []string
	ChapLanguageIETF []string
	ChapCountry      []string
}

type ChapProcess struct {
	ChapProcessCodecID uint
	ChapProcessPrivate []byte
	ChapProcessCommand []ChapProcessCommand
}

type ChapProcessCommand struct {
	ChapProcessTime uint
	ChapProcessData []byte
}

type Tags struct {
	Tag []Tag
}

type Tag struct {
	Targets   Targets
	SimpleTag []SimpleTag
}

type Targets struct {
	TargetTypeValue  uint
	TargetType       string
	TagTrackUID      []uint
	TagEditionUID    []uint
	TagChapterUID    []uint
	TagAttachmentUID []uint
}

type SimpleTag struct {
	TagName         string
	TagLanguage     string
	TagLanguageIETF string
	TagDefault      uint
	TagString       string
	TagBinary       []byte
}
