// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig(in *jlexer.Lexer, out *ELBServiceConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "CurrEnv":
			out.CurrEnv = string(in.String())
		case "ELBEnvConfigs":
			if in.IsNull() {
				in.Skip()
				out.ELBEnvConfigs = nil
			} else {
				in.Delim('[')
				if out.ELBEnvConfigs == nil {
					if !in.IsDelim(']') {
						out.ELBEnvConfigs = make([]ELBEnvConfig, 0, 1)
					} else {
						out.ELBEnvConfigs = []ELBEnvConfig{}
					}
				} else {
					out.ELBEnvConfigs = (out.ELBEnvConfigs)[:0]
				}
				for !in.IsDelim(']') {
					var v1 ELBEnvConfig
					(v1).UnmarshalEasyJSON(in)
					out.ELBEnvConfigs = append(out.ELBEnvConfigs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig(out *jwriter.Writer, in ELBServiceConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"CurrEnv\":"
		out.RawString(prefix[1:])
		out.String(string(in.CurrEnv))
	}
	{
		const prefix string = ",\"ELBEnvConfigs\":"
		out.RawString(prefix)
		if in.ELBEnvConfigs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.ELBEnvConfigs {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ELBServiceConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ELBServiceConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ELBServiceConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ELBServiceConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig(l, v)
}
func easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig1(in *jlexer.Lexer, out *ELBOpers) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "BindURL":
			out.BindURL = string(in.String())
		case "UnBindURL":
			out.UnBindURL = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig1(out *jwriter.Writer, in ELBOpers) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"BindURL\":"
		out.RawString(prefix[1:])
		out.String(string(in.BindURL))
	}
	{
		const prefix string = ",\"UnBindURL\":"
		out.RawString(prefix)
		out.String(string(in.UnBindURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ELBOpers) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ELBOpers) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ELBOpers) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ELBOpers) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig1(l, v)
}
func easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig2(in *jlexer.Lexer, out *ELBEnvConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Env":
			out.Env = string(in.String())
		case "ELBOpers":
			(out.ELBOpers).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig2(out *jwriter.Writer, in ELBEnvConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Env\":"
		out.RawString(prefix[1:])
		out.String(string(in.Env))
	}
	{
		const prefix string = ",\"ELBOpers\":"
		out.RawString(prefix)
		(in.ELBOpers).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ELBEnvConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ELBEnvConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeCalmwuOrgElbserviceOperatorPkgConfig2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ELBEnvConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ELBEnvConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeCalmwuOrgElbserviceOperatorPkgConfig2(l, v)
}