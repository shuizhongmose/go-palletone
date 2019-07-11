/*
 *
 *    This file is part of go-palletone.
 *    go-palletone is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU General Public License as published by
 *    the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *    go-palletone is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU General Public License for more details.
 *    You should have received a copy of the GNU General Public License
 *    along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
 * /
 *
 *  * @author PalletOne core developer <dev@pallet.one>
 *  * @date 2018-2019
 *
 */

package memunit

import (
	"fmt"
	"sync"
	"time"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/event"
	"github.com/palletone/go-palletone/common/hexutil"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/ptndb"
	"github.com/palletone/go-palletone/core"
	common2 "github.com/palletone/go-palletone/dag/common"
	"github.com/palletone/go-palletone/dag/errors"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/palletcache"
	"github.com/palletone/go-palletone/dag/txspool"
	"go.dedis.ch/kyber/v3/sign/bls"
)

type MemDag struct {
	token              modules.AssetId
	stableUnitHash     common.Hash
	stableUnitHeight   uint64
	lastMainChainUnit  *modules.Unit
	threshold          int
	orphanUnits        sync.Map
	orphanUnitsParants sync.Map
	chainUnits         sync.Map //Key:UnitHash,Value:UnitTempDb
	//leafUnits sync.Map
	//tempdbunitRep      map[common.Hash]common2.IUnitRepository
	//tempUtxoRep        map[common.Hash]common2.IUtxoRepository
	//tempStateRep       map[common.Hash]common2.IStateRepository
	//tempPropRep        map[common.Hash]common2.IPropRepository
	//tempUnitProduceRep common2.IUnitProduceRepository
	ldb               ptndb.Database
	ldbunitRep        common2.IUnitRepository
	ldbutxoRep        common2.IUtxoRepository
	ldbStateRep       common2.IStateRepository
	ldbPropRep        common2.IPropRepository
	ldbUnitProduceRep common2.IUnitProduceRepository
	//tempdb            map[common.Hash]*Tempdb
	saveHeaderOnly bool
	lock           sync.RWMutex
	//validator         validator.Validator
	cache palletcache.ICache
	// append by albert·gou 用于通知群签名
	toGroupSignFeed  event.Feed
	toGroupSignScope event.SubscriptionScope
}

func (pmg *MemDag) Close() {
	pmg.toGroupSignScope.Close()
}

func (pmg *MemDag) SubscribeToGroupSignEvent(ch chan<- modules.ToGroupSignEvent) event.Subscription {
	return pmg.toGroupSignScope.Track(pmg.toGroupSignFeed.Subscribe(ch))
}

func (pmg *MemDag) SetStableThreshold(count int) {
	pmg.lock.Lock()
	defer pmg.lock.Unlock()
	pmg.threshold = count
}

func NewMemDag(token modules.AssetId, threshold int, saveHeaderOnly bool, ldb ptndb.Database,
	stableUnitRep common2.IUnitRepository, stableUtxoRep common2.IUtxoRepository,
	stableStateRep common2.IStateRepository, propRep common2.IPropRepository, cache palletcache.ICache) *MemDag {

	ldbUnitProduceRep := common2.NewUnitProduceRepository(stableUnitRep, propRep, stableStateRep)
	stablehash, stbIndex, err := propRep.GetNewestUnit(token)
	if err != nil {
		log.Errorf("Cannot retrieve last stable unit from db for token:%s, you forget 'gptn init'??", token.String())
		return nil
	}
	var stableUnit *modules.Unit
	if saveHeaderOnly {
		header, err := stableUnitRep.GetHeaderByHash(stablehash)
		if err != nil {
			log.Errorf("Cannot retrieve last stable unit from db by hash[%s]", stablehash.String())
			return nil
		}
		stableUnit = modules.NewUnit(header, nil)
	} else {
		stableUnit, err = stableUnitRep.GetUnit(stablehash)
		if err != nil {
			log.Errorf("Cannot retrieve last stable unit from db by hash[%s]", stablehash.String())
			return nil
		}
	}
	log.Debugf("Init MemDag[%s], get last stable unit[%s] to set lastMainChainUnit", token.String(), stablehash.String())

	memdag := &MemDag{
		token:              token,
		threshold:          threshold,
		ldbunitRep:         stableUnitRep,
		ldbutxoRep:         stableUtxoRep,
		ldbPropRep:         propRep,
		ldbStateRep:        stableStateRep,
		ldbUnitProduceRep:  ldbUnitProduceRep,
		ldb:                ldb,
		orphanUnits:        sync.Map{},
		orphanUnitsParants: sync.Map{},
		chainUnits:         sync.Map{},
		stableUnitHash:     stablehash,
		stableUnitHeight:   stbIndex.Index,
		lastMainChainUnit:  stableUnit,
		saveHeaderOnly:     saveHeaderOnly,
		cache:              cache,
	}

	//
	//go memdag.loopRebuildTmpDb()
	return memdag
}

//func (chain *MemDag) loopRebuildTmpDb() {
//	rebuild := time.NewTicker(10 * time.Minute)
//	defer rebuild.Stop()
//	for {
//		select {
//		case <-rebuild.C:
//			if chain.lastMainChainUnit.Hash() == chain.stableUnitHash || len(chain.getChainUnits()) <= 1 {
//				// temp db don't need rebuild.
//				continue
//			}
//			chain.lock.Lock()
//			chain.rebuildTempdb()
//			chain.lock.Unlock()
//		}
//	}
//}
func (chain *MemDag) GetUnstableRepositories() (common2.IUnitRepository, common2.IUtxoRepository, common2.IStateRepository, common2.IPropRepository, common2.IUnitProduceRepository) {
	a, b, c, d, e, _ := chain.getUnstableRepositories()
	return a, b, c, d, e
}
func (chain *MemDag) getUnstableRepositories() (common2.IUnitRepository, common2.IUtxoRepository, common2.IStateRepository, common2.IPropRepository, common2.IUnitProduceRepository, error) {
	last_main_hash := chain.lastMainChainUnit.Hash()
	unittempdb, _ := chain.getChainUnit(last_main_hash)
	if unittempdb == nil { //没有任何不稳定单元
		return chain.ldbunitRep, chain.ldbutxoRep, chain.ldbStateRep, chain.ldbPropRep, chain.ldbUnitProduceRep, nil
	}
	return unittempdb.UnitRep, unittempdb.UtxoRep, unittempdb.StateRep, unittempdb.PropRep, unittempdb.UnitProduceRep, nil
}
func (chain *MemDag) GetHeaderByHash(hash common.Hash) (*modules.Header, error) {
	last_main_hash := chain.lastMainChainUnit.Hash()
	unittempdb, _ := chain.getChainUnit(last_main_hash)
	return unittempdb.UnitRep.GetHeaderByHash(hash)

}
func (chain *MemDag) GetHeaderByNumber(number *modules.ChainIndex) (*modules.Header, error) {
	last_main_hash := chain.lastMainChainUnit.Hash()
	unittempdb, _ := chain.getChainUnit(last_main_hash)
	return unittempdb.UnitRep.GetHeaderByNumber(number)
}

//func (chain *MemDag) getHeaderByNumber(number *modules.ChainIndex) (*modules.Header, error) {
//	chain_units := chain.getChainUnits()
//	for _, unit := range chain_units {
//		if index := unit.UnitHeader.ChainIndex(); index.Equal(number) {
//			return unit.Header(), nil
//		}
//	}
//	return nil, fmt.Errorf("the header[%s] not exist.", number.String())
//}

func (chain *MemDag) SetUnitGroupSign(uHash common.Hash /*, groupPubKey []byte*/, groupSign []byte,
	txpool txspool.ITxPool) error {
	//1. Set this unit as stable
	unit, err := chain.getChainUnit(uHash)
	if err != nil {
		log.Debugf("get Chain Unit error: %v", err.Error())
		return err
	}

	if !(unit.Unit.NumberU64() > chain.stableUnitHeight) {
		return nil
	}

	chain.lock.Lock()
	defer chain.lock.Unlock()
	chain.setStableUnit(uHash, unit.Unit.NumberU64(), txpool)

	//2. Update unit.groupSign
	header := unit.Unit.Header()
	//header.GroupPubKey = groupPubKey
	header.GroupSign = groupSign
	log.Debugf("Try to update unit[%s] header group sign", uHash.String())
	// 进行下一个unit的群签名
	if err == nil {
		log.Debugf("send toGroupSign event")
		go chain.toGroupSignFeed.Send(modules.ToGroupSignEvent{})
	}

	return chain.ldbunitRep.SaveHeader(header)
}

//设置某个单元和高度为稳定单元。设置后会更新当前的稳定单元，并将所有稳定单元写入到StableDB中，并且将ChainUnit中的稳定单元删除。
//然后基于最新的稳定单元，重建Tempdb数据库
func (chain *MemDag) setStableUnit(hash common.Hash, height uint64, txpool txspool.ITxPool) {
	tt := time.Now()
	log.Debugf("Set stable unit to %s,height:%d", hash.String(), height)
	stable_height := chain.stableUnitHeight
	stableCount := int(height - stable_height)
	if stableCount <= 0 {
		log.Errorf("Current stable height is %d, impossible set stable height to %d", stable_height, height)
		return
	}
	newStableUnits := make([]*modules.Unit, stableCount)
	stbHash := hash
	chain_units := chain.getChainUnits()
	for i := 0; i < stableCount; i++ {
		if u, has := chain_units[stbHash]; has {
			newStableUnits[stableCount-i-1] = u.Unit
			stbHash = u.Unit.ParentHash()[0]
		}
	}
	//Save stable unit and it's parent
	max_height := height
	for _, unit := range newStableUnits {
		if unit.NumberU64() > max_height {
			max_height = unit.NumberU64()
		}
		chain.setNextStableUnit(unit, txpool)
	}
	log.InfoDynamic(func() string {
		return fmt.Sprintf("set next stable unit cost time: %s ,index: %d, hash: %s",
			time.Since(tt), height, hash.String())
	})
	//remove fork units, and remove lower than stable unit
	for _, funit := range chain_units {
		if funit.Unit.NumberU64() <= max_height && funit.Unit.Hash() != hash {
			chain.removeUnitAndChildren(funit.Unit.Hash(), txpool)
		}
	}
	//remove too low orphan unit
	chain.removeLowOrphanUnit(max_height)
}

//设置当前稳定单元的指定父单元为稳定单元
func (chain *MemDag) setNextStableUnit(unit *modules.Unit, txpool txspool.ITxPool) {
	hash := unit.Hash()
	height := unit.NumberU64()
	unittempdb, _ := chain.getChainUnit(hash)

	// memdag不依赖apply unit的存储，因此用协程提高setStable的效率
	// 虽然与memdag无关，但是下一个unit的 apply 处理依赖上一个unit apply的结果，所以不能用协程并发处理
	unittempdb.Tempdb.PushData2ParentDb()
	if !chain.saveHeaderOnly && len(unit.Txs) > 1 {
		go txpool.SendStoredTxs(unit.Txs.GetTxIds())
	}
	chain.setChildrenUnitParentDb(hash, unittempdb.Tempdb.db)
	log.Debugf("Remove unit[%s] from chainUnits", hash.String())
	//remove new stable unit
	chain.chainUnits.Delete(hash)
	//Set stable unit
	chain.stableUnitHash = hash
	chain.stableUnitHeight = height
}
func (chain *MemDag) setChildrenUnitParentDb(parentHash common.Hash, db ptndb.Database) {
	for _, unit := range chain.getChainUnits() {
		if unit.Unit.ParentHash()[0] == parentHash {
			unit.SetParentDb(db)
		}
	}
}
func (chain *MemDag) checkUnitIrreversibleWithGroupSign(unit *modules.Unit) bool {
	if unit.GetGroupPubKeyByte() == nil || unit.GetGroupSign() == nil {
		return false
	}

	pubKey, err := unit.GetGroupPubKey()
	if err != nil {
		log.Debug(err.Error())
		return false
	}

	err = bls.Verify(core.Suite, pubKey, unit.UnitHash[:], unit.GetGroupSign())
	if err != nil {
		log.Debug(err.Error())
		return false
	}

	return true
}

// 判断当前主链上的单元是否有满足稳定单元的确认数，如果有，则更新稳定单元，并重建Temp数据库，返回True
// 如果没有，则不进行任何操作，返回False
func (chain *MemDag) checkStableCondition(unit *modules.Unit, txpool txspool.ITxPool) bool {
	// append by albert, 使用群签名判断是否稳定
	if chain.checkUnitIrreversibleWithGroupSign(unit) {
		log.Debugf("the unit(%s) have group sign(%s), make it to irreversible.",
			unit.UnitHash.TerminalString(), hexutil.Encode(unit.GetGroupSign()))
		chain.setStableUnit(unit.UnitHash, unit.NumberU64(), txpool)
		return true
	}

	unstableCount := int(unit.NumberU64() - chain.stableUnitHeight)
	//每个单元被多少个地址确认过(包括自己)
	unstableCofirmAddrs := make(map[common.Hash]map[common.Address]bool)
	childrenCofirmAddrs := make(map[common.Address]bool)
	ustbHash := unit.Hash()
	childrenCofirmAddrs[unit.Author()] = true
	units := chain.getChainUnits()
	// todo Albert·gou 待重做 优化逻辑
	for i := 0; i < unstableCount; i++ {
		u := units[ustbHash].Unit
		hs := unstableCofirmAddrs[ustbHash]
		if hs == nil {
			hs = make(map[common.Address]bool)
			unstableCofirmAddrs[ustbHash] = hs
		}
		hs[u.Author()] = true
		for addr := range childrenCofirmAddrs {
			hs[addr] = true
		}
		childrenCofirmAddrs[u.Author()] = true

		if len(hs) >= chain.threshold {
			log.Debugf("Unit[%s] has enough confirm address count=%d, make it to stable.", ustbHash.String(), len(hs))
			chain.setStableUnit(ustbHash, u.NumberU64(), txpool)
			return true
		}
		ustbHash = u.ParentHash()[0]
	}
	return false
}

//清空Tempdb，然后基于稳定单元到最新主链单元的路径，构建新的Tempdb
//func (chain *MemDag) rebuildTempdb() {
//	last_main_hash := chain.lastMainChainUnit.Hash()
//	forks := make([]common.Hash, 0)
//	for hash, temp := range chain.tempdb {
//		forks = append(forks, hash)
//		temp.Clear()
//	}
//
//	unstableUnits := chain.getMainChainUnits()
//	for _, unit := range unstableUnits {
//		chain.saveUnitToDb(chain.tempdbunitRep[last_main_hash], chain.tempUnitProduceRep, unit)
//	}
//}

//获得从稳定单元到最新单元的主链上的单元列表，从久到新排列
func (chain *MemDag) getMainChainUnits() []*modules.Unit {
	unstableCount := int(chain.lastMainChainUnit.NumberU64() - chain.stableUnitHeight)
	log.Debugf("Unstable unit count:%d", unstableCount)
	unstableUnits := make([]*modules.Unit, unstableCount)
	ustbHash := chain.lastMainChainUnit.Hash()
	chain_units := chain.getChainUnits()
	log.DebugDynamic(func() string {
		str := "chainUnits has unit:"
		for hash, _ := range chain_units {
			str += hash.String() + ";"
		}
		return str
	})
	for i := 0; i < unstableCount; i++ {
		u, ok := chain_units[ustbHash]
		if !ok {
			log.Errorf("chainUnits don't have unit[%s], last_main[%s]", ustbHash.String(), chain.lastMainChainUnit.Hash().String())
		}
		unstableUnits[unstableCount-i-1] = u.Unit
		ustbHash = u.Unit.ParentHash()[0]
	}
	return unstableUnits
}

func (chain *MemDag) getForkUnits(unit *modules.Unit) []*modules.Unit {
	chain_units := chain.getChainUnits()
	hash := unit.Hash()
	unstableCount := int(unit.NumberU64() - chain.stableUnitHeight)
	unstableUnits := make([]*modules.Unit, unstableCount)
	for i := 0; i < unstableCount; i++ {
		u, ok := chain_units[hash]
		if !ok {
			log.Errorf("getforks chainUnits don't have unit[%s], last_main[%s]", hash.String(), unit.Hash().String())
		}
		unstableUnits[unstableCount-i-1] = u.Unit
		hash = u.Unit.ParentHash()[0]
	}
	return unstableUnits
}

//判断当前设置是保存Header还是Unit，将对应的对象保存到Tempdb数据库
func (chain *MemDag) saveUnitToDb(unitRep common2.IUnitRepository, produceRep common2.IUnitProduceRepository, unit *modules.Unit) {
	log.Debugf("Save unit[%s] to db", unit.Hash().String())
	if chain.saveHeaderOnly {
		unitRep.SaveNewestHeader(unit.Header())
	} else {
		produceRep.PushUnit(unit)
	}
}

//从ChainUnits集合中删除一个单元以及其所有子孙单元
func (chain *MemDag) removeUnitAndChildren(hash common.Hash, txpool txspool.ITxPool) {
	log.Debugf("Remove unit[%s] and it's children from chain unit", hash.String())
	chain_units := chain.getChainUnits()
	for h, unit := range chain_units {
		if h == hash {
			if txs := unit.Unit.Transactions(); len(txs) > 1 {
				go txpool.ResetPendingTxs(txs)
			}
			chain.chainUnits.Delete(h)
			log.Debugf("Remove unit[%s] from chainUnits", hash.String())
		} else {
			if unit.Unit.ParentHash()[0] == hash {
				chain.removeUnitAndChildren(h, txpool)
			}
		}
	}
}

func (chain *MemDag) AddUnit(unit *modules.Unit, txpool txspool.ITxPool) (common2.IUnitRepository, common2.IUtxoRepository, common2.IStateRepository, common2.IPropRepository, common2.IUnitProduceRepository, error) {
	start := time.Now()
	if unit == nil {
		return nil, nil, nil, nil, nil, errors.ErrNullPoint
	}
	unitHash := unit.Hash()
	chain.lock.Lock()
	defer chain.lock.Unlock()
	if unit.NumberU64() <= chain.stableUnitHeight {
		log.Debugf("This unit is too old! Ignore it,stable unit height:%d, stable hash:%s",
			chain.stableUnitHeight, chain.stableUnitHash.String())
		return chain.getUnstableRepositories()
	}
	chain_units := chain.getChainUnits()
	if _, has := chain_units[unit.Hash()]; has { // 不重复添加
		log.Infof("MemDag[%s] received a repeated unit, hash[%s] ", chain.token.String(), unitHash.String())
		return chain.getUnstableRepositories()
	}
	err := chain.addUnit(unit, txpool)
	log.InfoDynamic(func() string {
		return fmt.Sprintf("MemDag[%s] AddUnit cost time: %v ,index: %d, hash: %s", chain.token.String(),
			time.Since(start), unit.NumberU64(), unitHash.String())
	})

	if err != nil {
		log.Warnf("Add unit[%s] fail, error:%s", unitHash.String(), err.Error())
		return nil, nil, nil, nil, nil, err
	} else {
		// 进行下一个unit的群签名
		log.Debugf("send toGroupSign event")
		go chain.toGroupSignFeed.Send(modules.ToGroupSignEvent{})
	}

	return chain.getUnstableRepositories()
}

func (chain *MemDag) addUnit(unit *modules.Unit, txpool txspool.ITxPool) error {
	parentHash := unit.ParentHash()[0]
	uHash := unit.Hash()
	height := unit.NumberU64()
	var parentDb ptndb.Database
	if parentHash == chain.stableUnitHash {
		parentDb = chain.ldb
	}
	if parentUnitTempdb, ok := chain.getChainUnits()[parentHash]; ok {
		parentDb = parentUnitTempdb.Tempdb
	}
	if parentDb != nil { //能找到父单元
		newTempDb, err := NewUnitTempDb(parentDb, unit, chain.cache, chain.saveHeaderOnly)
		if err != nil { //新来的单元无效
			return err
		}
		chain.chainUnits.Store(uHash, newTempDb)
		//add unit to chain
		log.Debugf("chain[%p] Add unit[%s] to chainUnits", chain, uHash.String())

		//add at the end of main chain unit
		if parentHash == chain.lastMainChainUnit.Hash() {
			//Add a new unit to main chain

			chain.setLastMainchainUnit(unit)

			//update txpool's tx status to pending
			if len(unit.Txs) > 0 {
				go txpool.SetPendingTxs(unit.Hash(), height, unit.Txs)
			}
			//增加了单元后检查是否满足稳定单元的条件
			start := time.Now()
			// todo Albert·gou 待重做 优化逻辑
			if chain.checkStableCondition(unit, txpool) {
				// 进行下一个unit的群签名
				log.Debugf("send toGroupSign event")
				go chain.toGroupSignFeed.Send(modules.ToGroupSignEvent{})
				log.Debugf("unit[%s] checkStableCondition =true", uHash.String())
			}
			log.InfoDynamic(func() string {
				return fmt.Sprintf("check stable cost time: %s ,index: %d, hash: %s",
					time.Since(start), height, uHash.String())
			})

			//start1 := time.Now()
			//unit_rep := chain.tempdbunitRep[parentHash]
			//chain.tempdbunitRep[uHash] = unit_rep
			//delete(chain.tempdbunitRep, parentHash)
			//chain.saveUnitToDb(unit_rep, chain.tempUnitProduceRep, unit)
			//log.DebugDynamic(func() string {
			//	return fmt.Sprintf("save unit cost time: %s ,index: %d, hash: %s",
			//		time.Since(start1), height, uHash.String())
			//})
		} else { //Fork unit

			// 满足切换主链条件， 则切换主链，更新主链单元。
			if unit.NumberU64() > chain.lastMainChainUnit.NumberU64() {
				//Need switch main chain
				chain.switchMainChain(unit, txpool)

				//log.InfoDynamic(func() string {
				//	return fmt.Sprintf("switch chain ,fork chain number:%d, main chain number:%d "+
				//		"index:%d ,hash:%s,main_hash:%s", height, chain.lastMainChainUnit.NumberU64(), uHash.String(), chain.lastMainChainUnit.Hash().String())
				//})
			}
		}
		//orphan unit can add below this unit?
		if inter, has := chain.orphanUnitsParants.Load(uHash); has {
			chain.orphanUnitsParants.Delete(uHash)
			next_hash := inter.(common.Hash)
			chain.processOrphan(next_hash, txpool)
		}
	} else { //找不到父单元，是孤儿单元
		//add unit to orphan
		log.Infof("This unit[%s] is an orphan unit", uHash.String())
		chain.orphanUnits.Store(uHash, unit)
		chain.orphanUnitsParants.Store(unit.ParentHash()[0], uHash)
	}
	return nil
}
func (chain *MemDag) getCofirmAddrs(hash common.Hash, height uint64) int {
	log.Infof("get confirm address index:%d,stable_height:%d,hash: %s", height, chain.stableUnitHeight, hash.String())
	num := 0
	count := int(height - chain.stableUnitHeight)
	unstableCofirmAddrs := make(map[common.Hash]map[common.Address]bool)
	childrenCofirmAddrs := make(map[common.Address]bool)
	chain_units := chain.getChainUnits()
	for i := 0; i < count; i++ {
		u := chain_units[hash].Unit
		hs := unstableCofirmAddrs[hash]
		if hs == nil {
			hs = make(map[common.Address]bool)
			unstableCofirmAddrs[hash] = hs
		}
		hs[u.Author()] = true
		for addr := range childrenCofirmAddrs {
			hs[addr] = true
		}
		childrenCofirmAddrs[u.Author()] = true

		hash = u.ParentHash()[0]
		num = len(hs)
	}
	return num
}

//计算一个单元到稳定单元之间有多少个确认地址数
func (chain *MemDag) getChainAddressCount(lastUnit *modules.Unit) int {
	addrs := map[common.Address]bool{}
	unitHash := lastUnit.Hash()
	units := chain.getChainUnits()
	for unitHash != chain.stableUnitHash {
		unit := units[unitHash].Unit
		addrs[unit.Author()] = true
		unitHash = unit.ParentHash()[0]
	}
	return len(addrs)
}

//发现一条更长的确认数更多的链，则放弃原有主链，切换成新主链
//1.将旧主链上包含的交易在交易池中重置
//2.设置新主链最新单元
//3.将新主链上包含的交易在交易池中设置为已打包
func (chain *MemDag) switchMainChain(newUnit *modules.Unit, txpool txspool.ITxPool) {
	oldMainchainUnits := chain.getMainChainUnits()

	// 不在主链上的区块，将已打包交易回滚。
	for _, un_unit := range oldMainchainUnits {
		txs := un_unit.Transactions()
		if len(txs) > 1 {
			// 用协程，resettPending和txpool的读写锁，会导致这里会很耗时。
			go txpool.ResetPendingTxs(txs)
		}
	}
	chain.setLastMainchainUnit(newUnit)
	newMainchainUnits := chain.getMainChainUnits()
	//基于新主链，更新TxPool的状态
	for _, unit := range newMainchainUnits {
		if len(unit.Txs) > 1 {
			log.Debugf("Update tx[%#x] status to pending in txpool", unit.Txs.GetTxIds())
			go txpool.SetPendingTxs(unit.Hash(), unit.NumberU64(), unit.Txs)
		}
	}
}

//将其从孤儿单元列表中删除，并添加到ChainUnits中。
func (chain *MemDag) processOrphan(unitHash common.Hash, txpool txspool.ITxPool) {
	unit, has := chain.getOrphanUnits()[unitHash]
	if has {
		log.Debugf("Orphan unit[%s] can add to chain now.", unit.Hash().String())
		chain.orphanUnits.Delete(unitHash)
		err := chain.addUnit(unit, txpool)
		if err != nil {
			log.Warnf("Add orphan unit[%s] to chain fail, error:%s", unitHash.String(), err.Error())
		}
	}
}
func (chain *MemDag) getOrphanUnits() map[common.Hash]*modules.Unit {
	units := make(map[common.Hash]*modules.Unit)
	chain.orphanUnits.Range(func(k, v interface{}) bool {
		hash := k.(common.Hash)
		u := v.(*modules.Unit)
		u_hash := u.Hash()
		if hash != u_hash {
			chain.orphanUnits.Delete(hash)
			chain.orphanUnits.Store(u_hash, u)
		}
		units[u_hash] = u
		return true
	})
	return units
}
func (chain *MemDag) removeLowOrphanUnit(lessThan uint64) {
	for hash, unit := range chain.getOrphanUnits() {
		if unit.NumberU64() <= lessThan {
			log.Debugf("Orphan unit[%s] height[%d] is too low, remove it.", unit.Hash().String(), unit.NumberU64())
			//if txs := unit.Transactions(); len(txs) > 1 {
			//	go txpool.ResetPendingTxs(txs)
			//}
			chain.orphanUnits.Delete(hash)
		}
	}
}
func (chain *MemDag) getChainUnits() map[common.Hash]*UnitTempDb {
	units := make(map[common.Hash]*UnitTempDb)
	chain.chainUnits.Range(func(k, v interface{}) bool {
		hash := k.(common.Hash)
		u := v.(*UnitTempDb)
		units[hash] = u
		return true
	})
	return units
}
func (chain *MemDag) getChainUnit(hash common.Hash) (*UnitTempDb, error) {
	inter, ok := chain.chainUnits.Load(hash)
	if ok {
		return inter.(*UnitTempDb), nil
	}
	return nil, errors.ErrNotFound
}

func (chain *MemDag) GetLastMainChainUnit() *modules.Unit {
	chain.lock.RLock()
	defer chain.lock.RUnlock()
	return chain.lastMainChainUnit
}

//设置最新的主链单元，并更新PropDB
func (chain *MemDag) setLastMainchainUnit(unit *modules.Unit) {
	// update tempdb interface
	//old_main_unit_hash := chain.lastMainChainUnit.Hash()
	//hash := unit.Hash()
	//chain.tempdbunitRep[hash] = chain.tempdbunitRep[old_main_unit_hash]
	//chain.tempPropRep[hash] = chain.tempPropRep[old_main_unit_hash]
	//chain.tempStateRep[hash] = chain.tempStateRep[old_main_unit_hash]
	//chain.tempUtxoRep[hash] = chain.tempUtxoRep[old_main_unit_hash]
	//chain.tempdb[hash] = chain.tempdb[old_main_unit_hash]
	//delete(chain.tempdbunitRep, old_main_unit_hash)
	//delete(chain.tempPropRep, old_main_unit_hash)
	//delete(chain.tempStateRep, old_main_unit_hash)
	//delete(chain.tempUtxoRep, old_main_unit_hash)
	//delete(chain.tempdb, old_main_unit_hash)

	chain.lastMainChainUnit = unit
}

//查询所有不稳定单元（不包括孤儿单元）
func (chain *MemDag) GetChainUnits() map[common.Hash]*modules.Unit {
	result := make(map[common.Hash]*modules.Unit)
	for uhash, db := range chain.getChainUnits() {
		result[uhash] = db.Unit
	}
	return result
}
