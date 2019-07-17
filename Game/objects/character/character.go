package character

import "time"
import _ "../ability"
import charmodel "../charmodel"
import "math/rand"

type Race byte

const (
	none Race = 0 + iota
	nuian
	fairy
	dwarf
	elf
	hariharan
	ferre
	returned
	warborn
)

type Gender byte

const (
	male Gender = 1 + iota
	female
)

// Character ...
type Character struct {

	//logger _log = LogManager.GetCurrentClassLogger();
	//options Dictionary<ushort, string> ;
	//GameConnection Connection net.Conn
	//List<IDisposable> Subscribers

	ID                  uint
	AccountID           uint64
	Race                byte
	Gender              byte
	LaborPower          uint16
	LaborPowerModified  time.Duration
	ConsumedLaborPower  int
	Ability1            byte
	Ability2            byte
	Ability3            byte
	FactionName         string
	Family              uint
	DeadCount           uint16
	DeadTime            time.Duration
	RezWaitDuration     int
	RezTime             time.Duration
	RezPenaltyDuration  int
	LeaveTime           time.Duration
	Money               int64
	Money2              int64
	HonorPoint          int
	VocationPoint       int
	CrimePoint          int
	CrimeRecord         int
	DeleteRequestTime   time.Duration
	TransferRequestTime time.Duration
	DeleteTime          time.Duration
	BmPoint             int64
	AutoUseAAPoint      bool
	PrevPoint           int
	Point               int
	Gift                int
	Expirience          int
	RecoverableExp      int
	Updated             time.Duration

	ReturnDictrictID       uint
	ResurrectionDictrictID uint

	/*
			override UnitCustomModelParams ModelParams
			override float Scale => 1f;
			override byte RaceGender => (byte)(16 * (byte)Gender + (byte)Race);

			CharacterVisualOptions VisualOptions

			ActionSlot[] Slots
			Inventory Inventory
			byte NumInventorySlots
			short NumBankSlots

			Item[] BuyBack
			BondDoodad Bonding
			CharacterQuests Quests
			CharacterMails Mails
			CharacterAppellations Appellations
			CharacterAbilities Abilities
			CharacterPortals Portals
			CharacterFriends Friends
			CharacterBlocked Blocked
			CharacterMates Mates

			byte ExpandedExpert
			CharacterActability Actability

			CharacterSkills Skills
			CharacterCraft Craft

			int AccessLevel

			private bool _inParty;
			private bool _isOnline;

			bool InParty
			{
				get => _inParty;
				set
				{
					if (_inParty == value) return;
					// TODO - GUILD STATUS CHANGE
					FriendMananger.Instance.SendStatusChange(this, false, value);
					_inParty = value;
				}
			}

			bool IsOnline
			{
				get => _isOnline;
				set
				{
					if (_isOnline == value) return;
					// TODO - GUILD STATUS CHANGE
					FriendMananger.Instance.SendStatusChange(this, true, value);
					if (!value) TeamManager.Instance.SetOffline(this);
					_isOnline = value;
				}
			}
		}
	*/
}

// New ... Returns object of character
func New() *Character {
	return &Character{}
}

// Create ... Validates character and saves it in DB
func (char *Character) Create(name string, race, gender, ability1 byte, body []uint32, customModel *charmodel.CharacterModel) {
	// TODO: Validate name, Generate ID
	// var characterId = CharacterIdManager.Instance.GetNextId();
	char.ID = uint(rand.Int())
	// NameManager.Instance.AddCharacterName(characterId, name);
	/*
		var template = GetTemplate(race, gender);
						var character = new Character(customModel);
						character.Id = characterId;
						character.AccountId = connection.AccountId;
						character.Name = name.Substring(0, 1).ToUpper() + name.Substring(1);
						character.Race = (Race)race;
						character.Gender = (Gender)gender;



						character.Position = template.Position.Clone();
						character.Position.ZoneId = template.ZoneId;
						character.Level = 1;
						character.Faction = FactionManager.Instance.GetFaction(template.FactionId);
						character.FactionName = "";
						character.LaborPower = 50;
						character.LaborPowerModified = DateTime.UtcNow;
						character.NumInventorySlots = template.NumInventorySlot;
						character.NumBankSlots = template.NumBankSlot;
						character.Inventory = new Inventory(character);
						character.Updated = DateTime.UtcNow;
						character.Ability1 = (AbilityType)ability1;
						character.Ability2 = AbilityType.None;
						character.Ability3 = AbilityType.None;
						character.ReturnDictrictId = template.ReturnDictrictId;
						character.ResurrectionDictrictId = template.ResurrectionDictrictId;
						character.Slots = new ActionSlot[85];
			/*
						for (var i = 0; i < character.Slots.Length; i++)
						{
							character.Slots[i] = new ActionSlot();
						}
						var items = _abilityItems[ability1];
						SetEquipItemTemplate(character.Inventory, items.Items.Headgear, EquipmentItemSlot.Head, items.Items.HeadgearGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Necklace, EquipmentItemSlot.Neck, items.Items.NecklaceGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Shirt, EquipmentItemSlot.Chest, items.Items.ShirtGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Belt, EquipmentItemSlot.Waist, items.Items.BeltGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Pants, EquipmentItemSlot.Legs, items.Items.PantsGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Gloves, EquipmentItemSlot.Hands, items.Items.GlovesGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Shoes, EquipmentItemSlot.Feet, items.Items.ShoesGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Bracelet, EquipmentItemSlot.Arms, items.Items.BraceletGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Back, EquipmentItemSlot.Back, items.Items.BackGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Undershirts, EquipmentItemSlot.Undershirt, items.Items.UndershirtsGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Underpants, EquipmentItemSlot.Underpants, items.Items.UnderpantsGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Mainhand, EquipmentItemSlot.Mainhand, items.Items.MainhandGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Offhand, EquipmentItemSlot.Offhand, items.Items.OffhandGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Ranged, EquipmentItemSlot.Ranged, items.Items.RangedGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Musical, EquipmentItemSlot.Musical, items.Items.MusicalGrade);
						SetEquipItemTemplate(character.Inventory, items.Items.Cosplay, EquipmentItemSlot.Cosplay, items.Items.CosplayGrade);
						for (var i = 0; i < 7; i++)
						{
							if (body[i] == 0 && template.Items[i] > 0)
							{
								body[i] = template.Items[i];
							}
							SetEquipItemTemplate(character.Inventory, body[i], (EquipmentItemSlot)(i + 19), 0);
						}

						byte slot = 10;
						foreach (var item in items.Supplies)
						{
							var createdItem = ItemManager.Instance.Create(item.Id, item.Amount, item.Grade);
							character.Inventory.AddItem(createdItem);

							character.SetAction(slot, ActionSlotType.Item, item.Id);
							slot++;
						}

						items = _abilityItems[0];
						if (items != null)
						{
							foreach (var item in items.Supplies)
							{
								var createdItem = ItemManager.Instance.Create(item.Id, item.Amount, item.Grade);
								character.Inventory.AddItem(createdItem);

								character.SetAction(slot, ActionSlotType.Item, item.Id);
								slot++;
							}
						}

						character.Abilities = new CharacterAbilities(character);
						character.Abilities.SetAbility(character.Ability1, 0);

						character.Actability = new CharacterActability(character);
						foreach (var (id, actabilityTemplate) in _actabilities)
						{
							character.Actability.Actabilities.Add(id, new Actability(actabilityTemplate));
						}

						character.Skills = new CharacterSkills(character);
						foreach (var skill in SkillManager.Instance.GetDefaultSkills())
						{
							if (!skill.AddToSlot)
							{
								continue;
							}
							character.SetAction(skill.Slot, ActionSlotType.Skill, skill.Template.Id);
						}

						slot = 1;
						while (character.Slots[slot].Type != ActionSlotType.None)
						{
							slot++;
						}
						foreach (var skill in SkillManager.Instance.GetStartAbilitySkills(character.Ability1))
						{
							character.Skills.AddSkill(skill, 1, false);
							character.SetAction(slot, ActionSlotType.Skill, skill.Id);
							slot++;
						}

						character.Appellations = new CharacterAppellations(character);
						character.Quests = new CharacterQuests(character);
						character.Mails = new CharacterMails(character);
						character.Portals = new CharacterPortals(character);
						character.Friends = new CharacterFriends(character);

						character.Hp = character.MaxHp;
						character.Mp = character.MaxMp;

						if (character.Save())
						{
							connection.Characters.Add(character.Id, character);
							connection.SendPacket(new SCCreateCharacterResponsePacket(character));
						}
						else
						{
							connection.SendPacket(new SCCharacterCreationFailedPacket(3));
							CharacterIdManager.Instance.ReleaseId(characterId);
							NameManager.Instance.RemoveCharacterName(characterId);
							// TODO release items...
						}
					}
					else
					{
						connection.SendPacket(new SCCharacterCreationFailedPacket(nameValidationCode));
					}

	*/
}
