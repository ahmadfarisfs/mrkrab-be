package store

import (
	"log"

	"github.com/ahmadfarisfs/krab-core/contract"
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type MutationStore struct {
	db *gorm.DB
	ps contract.ProjectStore
}

func NewMutationStore(projectStore contract.ProjectStore, db *gorm.DB) *MutationStore {
	return &MutationStore{
		db: db,
		ps: projectStore,
	}
}

type MutationProjectResult struct {
	model.Mutation
	model.Project
}
type MutationPocketResult struct {
	model.Budget
	model.Account
}

func (ms *MutationStore) ListMutation(req utils.CommonRequest) ([]model.MutationExtended, int, error) {

	if req.Filter["projectIDs"] != nil {
		log.Println("req.Filter[projectIDs]----------------->")

		log.Println(req.Filter["projectIDs"])
		if req.Filter["pocketIDs"] != nil {
			//complete filter
			log.Println("req.Filter[pocketIDs]----------------->")

			log.Println(req.Filter["pocketIDs"])
		} else {
			//project with all pocket under it
		}
	} else {
		//all projects and all pockets
		//treated as error, yakali ngapain

	}
	// if req.Filter["projectIDs"] == nil { //} || req.Filter["account_ids"].(type) != int[] {
	// 	return []model.MutationExtended{}, 0, errors.New("Invalid params: projectIDs")
	// }
	initQuery := ms.db.
		Joins("left join accounts as a on a.id = mutations.account_id").
		Joins("left JOIN accounts as a2  ON (a2.id = a.parent_id and a.parent_id is not null)").
		Joins("left join projects as prj on (prj.account_id=a.id and a.parent_id is null)").
		Joins("left join projects as prj2 on (a2.id = prj2.account_id)").
		Joins("left join budgets on (budgets.account_id=a.id and a.parent_id is not null)"). //akun yang parent id nya tidak null adalah budgets
		Joins("left join transactions as t on t.id = mutations.transaction_id").
		Where("prj.id is not null or prj2.id is not null").
		Order("mutations.created_at desc").Table("mutations").
		Select("transactions.transaction_date", "mutations.id", " budgets.id as pocket_id", "coalesce(prj.id,prj2.id) as project_id", "coalesce(prj.is_open,prj2.is_open ) as is_open", "mutations.created_at", "mutations.amount ", "t.remarks ", "t.transaction_code ", "coalesce (prj.description,prj2.description )as project_description", "coalesce (prj.name,prj2.name ) as project_name", "budgets.name as pocket_name ", "budgets.`limit` as pocket_limit")

	//find in tabel akun (project)
	// projectDetails := []model.Project{}
	// queryProj := ms.db.Model(&model.Project{})
	if req.Filter["projectIDs"] != nil {
		initQuery = initQuery.Where("(prj2.id in (?) or prj.id in (?))", req.Filter["projectIDs"], req.Filter["projectIDs"])
		delete(req.Filter, "projectIDs")
	}
	// queryProj.Find(&projectDetails)

	//find in tabel akun (pocket)
	// pocketDetails := []model.Budget{}
	// queryPocket := ms.db.Model(&model.Budget{})
	if req.Filter["pocketIDs"] != nil {
		// initQuery = initQuery.Where("(budgets.id is null or budgets.id in (?))", req.Filter["pocketIDs"])
		initQuery = initQuery.Where(" budgets.id in (?)", req.Filter["pocketIDs"])

		delete(req.Filter, "pocketIDs")
	}

	if req.Filter["type"] != nil {
		//filter by expense,all, or income
		switch req.Filter["type"] {
		case "expense":
			initQuery = initQuery.Where("mutations.amount < 0")
		case "income":
			initQuery = initQuery.Where("mutations.amount > 0")
		}
		delete(req.Filter, "type")
	}
	if req.Filter["start_date"] != nil {
		initQuery = initQuery.Where("mutations.created_at > ?", req.Filter["start_date"])
		delete(req.Filter, "start_date")
	}
	if req.Filter["end_date"] != nil {
		initQuery = initQuery.Where("mutations.created_at < ?", req.Filter["end_date"])

		delete(req.Filter, "end_date")
	}

	// initQuery.Find(&pocketDetails)
	//===========================================

	// projectAccountIDs := []uint{}
	// pocketAccountIDs := []uint{}
	// allAccountIDs := []uint{}
	// for _, e := range projectDetails {
	// 	projectAccountIDs = append(projectAccountIDs, e.ID)
	// 	allAccountIDs = append(allAccountIDs, e.ID)
	// }
	// for _, e := range pocketDetails {
	// 	pocketAccountIDs = append(pocketAccountIDs, e.ID)
	// 	allAccountIDs = append(allAccountIDs, e.ID)
	// }
	// rawQuery:=`select m.id, bgt.id as pocket_id,coalesce (prj.id,prj2.id) as project_id,coalesce (prj.is_open,prj2.is_open ) as is_open,m.created_at,m.amount ,t.remarks ,t.transaction_code ,coalesce (prj.description,prj2.description )as project_description,coalesce (prj.name,prj2.name ) as project_name,bgt.name as pocket_name ,bgt.``limit`` as pocket_limit from mrkrab.mutations m
	// left join mrkrab.accounts a on a.id = m.account_id
	// left JOIN mrkrab.accounts a2  ON (a2.id = a.parent_id and a.parent_id is not null)
	// left join mrkrab.projects prj on (prj.account_id=a.id and a.parent_id is null)
	// left join mrkrab.projects prj2 on (a2.id = prj2.account_id)
	// left join mrkrab.budgets bgt on (bgt.account_id=a.id and a.parent_id is not null)
	// left join mrkrab.transactions t on t.id = m.transaction_id
	// where (prj.id is not null or prj2.id is not null)
	// order by created_at asc`

	// projectAccount := []MutationProjectResult{}
	// ms.db.Where("account_id IN(?)", projectAccountIDs).Model(&model.Mutation{}).
	// 	Joins("join projects ON projects.account_id = mutations.account_id").
	// 	Scan(&projectAccount)
	// log.Println(json.Marshal(projectAccount))

	// pocketAccount := []MutationPocketResult{}
	// ms.db.Where("account_id IN(?)", pocketAccountIDs).Model(&model.Mutation{}).
	// 	Joins("join budgets ON budgets.account_id = mutations.account_id").
	// 	Scan(&pocketAccount)
	// log.Println(json.Marshal(pocketAccount))

	// ms.db.Where()

	ret := []model.MutationExtended{}
	var count int64
	//query builder
	// initQuery := ms.db.Where("account_id IN(?)", projectAccountIDs)

	err := initQuery.Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	log.Println(req)
	//actually fetch data with limit and offset, only
	quer := utils.AppendCommonRequest(initQuery, req)
	err = quer.Find(&ret).Error

	return ret, int(count), err

	//	return nil, 0, nil
}
