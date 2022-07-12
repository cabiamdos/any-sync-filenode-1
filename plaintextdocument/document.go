package plaintextdocument

import (
	"fmt"

	"github.com/anytypeio/go-anytype-infrastructure-experiments/account"
	aclpb "github.com/anytypeio/go-anytype-infrastructure-experiments/aclchanges/pb"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/acltree"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/testutils/testchanges/pb"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/thread"
	"github.com/gogo/protobuf/proto"
)

type PlainTextDocument interface {
	Text() string
	AddText(text string) error
}

type plainTextDocument struct {
	heads   []string
	aclTree acltree.ACLTree
	state   *DocumentState
}

func (p *plainTextDocument) Text() string {
	return p.state.Text
}

func (p *plainTextDocument) AddText(text string) error {
	_, err := p.aclTree.AddContent(func(builder acltree.ChangeBuilder) error {
		builder.AddChangeContent(
			&pb.PlainTextChangeData{
				Content: []*pb.PlainTextChangeContent{
					createAppendTextChangeContent(text),
				},
			})
		return nil
	})
	return err
}

func (p *plainTextDocument) Update(tree acltree.ACLTree) {
	p.aclTree = tree
	var err error
	defer func() {
		if err != nil {
			fmt.Println("rebuild has returned error:", err)
		}
	}()

	prevHeads := p.heads
	p.heads = tree.Heads()
	startId := prevHeads[0]
	tree.IterateFrom(startId, func(change *acltree.Change) (isContinue bool) {
		if change.Id == startId {
			return true
		}
		if change.DecryptedDocumentChange != nil {
			p.state, err = p.state.ApplyChange(change.DecryptedDocumentChange, change.Id)
			if err != nil {
				return false
			}
		}
		return true
	})
}

func (p *plainTextDocument) Rebuild(tree acltree.ACLTree) {
	p.aclTree = tree
	p.heads = tree.Heads()
	var startId string
	var err error
	defer func() {
		if err != nil {
			fmt.Println("rebuild has returned error:", err)
		}
	}()

	rootChange := tree.Root()

	if rootChange.DecryptedDocumentChange == nil {
		err = fmt.Errorf("root doesn't have decrypted change")
		return
	}

	state, err := BuildDocumentStateFromChange(rootChange.DecryptedDocumentChange, rootChange.Id)
	if err != nil {
		return
	}

	startId = rootChange.Id
	tree.Iterate(func(change *acltree.Change) (isContinue bool) {
		if startId == change.Id {
			return true
		}
		if change.DecryptedDocumentChange != nil {
			state, err = state.ApplyChange(change.DecryptedDocumentChange, change.Id)
			if err != nil {
				return false
			}
		}
		return true
	})
	if err != nil {
		return
	}
	p.state = state
}

func NewInMemoryPlainTextDocument(acc *account.AccountData, text string) (PlainTextDocument, error) {
	return NewPlainTextDocument(acc, thread.NewInMemoryThread, text)
}

func NewPlainTextDocument(
	acc *account.AccountData,
	create func(change *thread.RawChange) (thread.Thread, error),
	text string) (PlainTextDocument, error) {
	changeBuilder := func(builder acltree.ChangeBuilder) error {
		err := builder.UserAdd(acc.Identity, acc.EncKey.GetPublic(), aclpb.ACLChange_Admin)
		if err != nil {
			return err
		}
		builder.AddChangeContent(createInitialChangeContent(text))
		return nil
	}
	t, err := acltree.BuildThreadWithACL(
		acc,
		changeBuilder,
		create)
	if err != nil {
		return nil, err
	}

	doc := &plainTextDocument{
		heads:   nil,
		aclTree: nil,
		state:   nil,
	}
	tree, err := acltree.BuildACLTree(t, acc, doc)
	if err != nil {
		return nil, err
	}
	doc.aclTree = tree
	return doc, nil
}

func createInitialChangeContent(text string) proto.Marshaler {
	return &pb.PlainTextChangeData{
		Content: []*pb.PlainTextChangeContent{
			createAppendTextChangeContent(text),
		},
		Snapshot: &pb.PlainTextChangeSnapshot{Text: text},
	}
}

func createAppendTextChangeContent(text string) *pb.PlainTextChangeContent {
	return &pb.PlainTextChangeContent{
		Value: &pb.PlainTextChangeContentValueOfTextAppend{
			TextAppend: &pb.PlainTextChangeTextAppend{
				Text: text,
			},
		},
	}
}
