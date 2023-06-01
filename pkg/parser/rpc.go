/*
Get the message definitions by parsing a proto file
No go files
*/
package parser

import (
	"fmt"
	"os"

	"github.com/emicklei/proto"
)

type protoParser struct {
	sequenceMap map[string]map[string]int
}

func (p *protoParser) GetSequence(messageName, fieldName string) int {
	if _, ok := p.sequenceMap[messageName]; !ok {
		p.sequenceMap[messageName] = map[string]int{}
		p.sequenceMap[messageName][fieldName] = 1
	}
	if _, ok := p.sequenceMap[messageName][fieldName]; !ok {
		// find the max sequence
		max := 0
		for _, v := range p.sequenceMap[messageName] {
			if v > max {
				max = v
			}
		}
		p.sequenceMap[messageName][fieldName] = max + 1
	}
	return p.sequenceMap[messageName][fieldName]
}

func ParseProtoFile(path string) (*protoParser, error) {
	reader, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &protoParser{sequenceMap: map[string]map[string]int{}}, nil
		}
		return nil, err
	}
	defer reader.Close()
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	p := &protoParser{
		sequenceMap: make(map[string]map[string]int),
	}

	proto.Walk(definition,
		proto.WithMessage(p.handleMessage),
		// proto.WithNormalField(VisitNormalField),
	)
	return p, nil
}

func (p *protoParser) handleMessage(m *proto.Message) {
	for _, each := range m.Elements {
		if f, ok := each.(*proto.NormalField); ok {
			if _, ok := p.sequenceMap[m.Name]; !ok {
				p.sequenceMap[m.Name] = make(map[string]int)
			}
			p.sequenceMap[m.Name][f.Name] = f.Sequence
		}
	}
}

type optionLister struct {
	proto.NoopVisitor
}

func (l optionLister) VisitOption(o *proto.Option) {
	fmt.Println(o.Name)
}

func (l optionLister) VisitMapField(m *proto.MapField) {
	fmt.Println(m.Name)
}

func VisitNormalField(m *proto.NormalField) {
	m.Accept(new(optionLister))
	fmt.Println(m.Name)
}
