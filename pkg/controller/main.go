package controller

import (
	"bytes"
	"fmt"
	"image"
	// Needed for decoding gif images with image
	_ "image/gif"
	"image/jpeg"
	// Needed for decoding png images with image
	_ "image/png"
	"log"
	"os"
	"os/user"
	"strconv"
	"time"
	"unicode"

	"github.com/Japodrilo/MyP-Proyecto2/pkg/model"
	"github.com/Japodrilo/MyP-Proyecto2/pkg/view"

	"github.com/dhowden/tag"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Principal is the main window controller. It contains as fields
// a database from the model package, a MainWindow object from the
// view package, a tree view, and the tree selection of the former.
type Principal struct {
	database   *model.Database
	mainWindow *view.MainWindow
	treeview   *TreeView
	treeSel    *gtk.TreeSelection
}

// A SongInfo holds the information of a Rola to show in the bottom
// of the main window, including the image obtained from the ID3v2 tag.
type SongInfo struct {
	image  *gtk.Image
	title  string
	artist string
	album  string
}

// MainWindow creates and draws the main window for the application.
// This is the function used in the 'main' package.
func MainWindow() {
	principal := newPrincipal()
	principal.initialize()
	principal.mainWindow.Win.ShowAll()
}

func newPrincipal() *Principal {
	database, exists := model.NewDatabase()
	mainWindow := view.SetupMainWindow()
	treeview := NewTreeView(mainWindow.TreeView)
	sel, err := treeview.TreeView.TreeView.GetSelection()
	if err != nil {
		log.Fatal("could not retrieve the treeview selection:", err)
	}

	if !exists {
		database.CreateDB()
	}

	principal := &Principal{
		database:   database,
		mainWindow: mainWindow,
		treeview:   treeview,
		treeSel:    sel,
	}
	return principal
}

func (principal *Principal) initialize() {
	principal.database.LoadDB()

	principal.mainWindow.Buttons["new"].Connect("clicked", func() {
		principal.addNewPerformer()
	})

	principal.mainWindow.Buttons["performers"].Connect("clicked", func() {
		principal.editForeignPerformer()
	})

	principal.treeSel.Connect("changed", func() {
		principal.selectionChanged(principal.treeSel)
	})

	principal.treeview.TreeView.TreeView.Connect("row-activated", func() {
		principal.rowActivated()
	})

	principal.mainWindow.Buttons["populate"].Connect("clicked", func() {
		principal.treeSel.UnselectAll()
		principal.treeSel.SetMode(gtk.SELECTION_NONE)
		principal.repopulate()
		principal.mainWindow.Buttons["populate"].SetSensitive(false)
		time.Sleep(100 * time.Nanosecond)
		principal.populate()
	})

	principal.mainWindow.Buttons["edit"].Connect("clicked", func() {
		principal.editPerformer()
	})

	principal.mainWindow.SearchEntry.Connect("activate", func() {
		text := view.GetTextEntryClean(principal.mainWindow.SearchEntry)
		principal.searchAction(text)
	})

	principal.mainWindow.Win.ShowAll()
}

func (principal *Principal) addNewPerformer() {
	performerPopUp := view.EditPerformerWindow()
	performerPopUp.Win.SetTitle("New Performer")
	performerPopUp.SaveB.Connect("clicked", func() {
		page := performerPopUp.Notebook.GetCurrentPage()
		switch page {
		case 0:
			principal.savePersonContent(performerPopUp.PersonContent)
		case 1:
			principal.saveGroupContent(performerPopUp.GroupContent)
		}
		performerPopUp.Win.Close()
	})
}

func (principal *Principal) populate() {
	miner := model.NewMiner()
	miner.Traverse()
	go miner.Extract()
	time.Sleep(100 * time.Millisecond)
	go miner.Populate(principal.database)
	time.Sleep(100 * time.Millisecond)
	go principal.populateOnTheFly(miner)
}

func (principal *Principal) repopulate() {
	principal.database.LoadDB()
	principal.populateFromExistingDB(principal.database)
}

func (principal *Principal) rowTextValues() []string {
	values := make([]string, 0)
	_, iter, ok := principal.treeSel.GetSelected()
	if !ok {
		return values
	}
	for i := 0; i < 5; i++ {
		cell, _ := principal.treeview.Filter.GetValue(iter, i)
		str, _ := cell.GetString()
		values = append(values, str)
	}
	return values
}

func (principal *Principal) rowID() int64 {
	_, iter, ok := principal.treeSel.GetSelected()
	if !ok {
		return -1
	}
	cell, _ := principal.treeview.Filter.GetValue(iter, 6)
	idU, _ := cell.GoValue()
	id := idU.(int)
	id64 := int64(id)
	return id64
}

// Handler of "activate" signal of TreeView's selection
func (principal *Principal) selectionChanged(s *gtk.TreeSelection) {
	home, err := user.Current()
	if err != nil {
		log.Fatal("could not retrieve the current user:", err)
	}
	cache := home.HomeDir + "/.cache/rolas"

	items := principal.rowTextValues()
	if len(items) == 0 {
		return
	}
	for i := 0; i < 4; i++ {
		if i < 3 {
			fmt.Print(items[i] + " - ")
		} else {
			fmt.Println(items[i])
		}
	}
	file, err := os.Open(items[4])
	if err != nil {
		log.Fatal("could not open file:", err)
	}
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal("error while reading the tags in file: ", items[4]+" ", err)
	}
	picture := metadata.Picture()
	if picture == nil {
		principal.defaultImage(items[0], items[1], items[2])
	} else {
		pic := picture.Data
		file, _ := os.Create(cache + "/image.jpg")
		imageReader := bytes.NewReader(pic)
		loadedImage, _, _ := image.Decode(imageReader)
		if loadedImage != nil {
			err = jpeg.Encode(file, loadedImage, nil)
			pix, _ := gdk.PixbufNewFromFileAtScale(cache+"/image.jpg", 250, 250, false)
			image, _ := gtk.ImageNewFromPixbuf(pix)
			glib.IdleAdd(principal.attachInfo, &SongInfo{image, items[0], items[1], items[2]})
			if err != nil {
				log.Fatal("could not encode the image to jpeg")
			}
		} else {
			principal.defaultImage(items[0], items[1], items[2])
		}
	}
}

func (principal *Principal) editPerformer() {
	rowValues := principal.rowTextValues()
	if len(rowValues) == 0 {
		return
	}
	rolaID := principal.rowID()
	rola := principal.database.QueryRola(rolaID)

	performerID := principal.database.ExistsPerformer(rola.Artist())
	ptype, name := principal.database.QueryPerformerType(performerID)
	switch ptype {
	case 0:
		personPopUp := view.EditPersonWindow()
		principal.showPersonContent(personPopUp.PersonContent, name)
		personID := principal.database.ExistsPerson(name)
		var listBoxRow *gtk.ListBoxRow
		personGroups := principal.database.QueryPersonGroups(personID)
		for group := range principal.database.AllGroups() {
			if !personGroups[group] {
				personPopUp.NewGroupCBT.AppendText(group)
			} else {
				listBoxRow = view.SetupListBoxRowLabel(group)
				listBoxRow.SetSensitive(false)
				personPopUp.CurrentGroupLB.Add(listBoxRow)
				personPopUp.Win.ShowAll()
			}
		}
		personPopUp.Notebook.ConnectAfter("switch-page", func() {
			page := personPopUp.Notebook.GetCurrentPage()
			switch page {
			case 0:
				glib.IdleAdd(personPopUp.SaveB.SetLabel, "Save")
			default:
				glib.IdleAdd(personPopUp.SaveB.SetLabel, "Add")
			}
		})
		personPopUp.SaveB.Connect("clicked", func() {
			page := personPopUp.Notebook.GetCurrentPage()
			switch page {
			case 0:
				principal.savePersonContent(personPopUp.PersonContent)
			case 1:
				principal.database.AddPersonToGroup(personID, principal.database.AllGroups()[personPopUp.NewGroupCBT.GetActiveText()])
			}
			personPopUp.Win.Close()
		})
	case 1:
		groupPopUp := view.EditGroupWindow()
		principal.showGroupContent(groupPopUp.GroupContent, name)
		groupID := principal.database.ExistsGroup(name)
		var listBoxRow *gtk.ListBoxRow
		groupMembers := principal.database.QueryGroupMembers(groupID)
		for member := range principal.database.AllPersons() {
			if !groupMembers[member] {
				groupPopUp.NewMemberCBT.AppendText(member)
			} else {
				listBoxRow = view.SetupListBoxRowLabel(member)
				listBoxRow.SetSensitive(false)
				groupPopUp.CurrentMemberLB.Add(listBoxRow)
				groupPopUp.Win.ShowAll()
			}
		}
		groupPopUp.Notebook.ConnectAfter("switch-page", func() {
			page := groupPopUp.Notebook.GetCurrentPage()
			switch page {
			case 0:
				glib.IdleAdd(groupPopUp.SaveB.SetLabel, "Save")
			default:
				glib.IdleAdd(groupPopUp.SaveB.SetLabel, "Add")
			}
		})
		groupPopUp.SaveB.Connect("clicked", func() {
			page := groupPopUp.Notebook.GetCurrentPage()
			switch page {
			case 0:
				principal.saveGroupContent(groupPopUp.GroupContent)
			case 1:
				principal.database.AddPersonToGroup(principal.database.AllPersons()[groupPopUp.NewMemberCBT.GetActiveText()], groupID)
			}
			groupPopUp.Win.Close()
		})
	case 2:
		performerPopUp := view.EditPerformerWindow()
		rola := principal.database.QueryRola(rolaID)
		performerPopUp.PersonContent.StageNameE.SetText(rola.Artist())
		performerPopUp.PersonContent.StageNameE.SetSensitive(false)
		performerPopUp.GroupContent.GroupNameE.SetText(rola.Artist())
		performerPopUp.GroupContent.GroupNameE.SetSensitive(false)
		performerPopUp.SaveB.Connect("clicked", func() {
			page := performerPopUp.Notebook.GetCurrentPage()
			switch page {
			case 0:
				principal.database.UpdatePerformerType(performerID, 0)
				principal.savePersonContent(performerPopUp.PersonContent)
			case 1:
				principal.database.UpdatePerformerType(performerID, 1)
				principal.saveGroupContent(performerPopUp.GroupContent)
			}
			performerPopUp.Win.Close()
		})
	}
}

func (principal *Principal) editForeignPerformer() {
	var groupID int64
	var groupName string
	var personID int64
	var personName string
	foreignPopUp := view.EditForeignPerformerWindow()
	for person := range principal.database.AllPersons() {
		foreignPopUp.PersonCBT.AppendText(person)
	}
	for group := range principal.database.AllGroups() {
		foreignPopUp.GroupCBT.AppendText(group)
	}
	foreignPopUp.PersonCBT.Connect("changed", func() {
		personName = foreignPopUp.PersonCBT.GetActiveText()
		principal.showPersonContent(foreignPopUp.PersonContent, personName)
		personID = principal.database.ExistsPerson(personName)
	})
	foreignPopUp.GroupCBT.Connect("changed", func() {
		groupName = foreignPopUp.GroupCBT.GetActiveText()
		principal.showGroupContent(foreignPopUp.GroupContent, groupName)
		groupID = principal.database.ExistsGroup(groupName)
	})
	foreignPopUp.SaveB.Connect("clicked", func() {
		switch foreignPopUp.Notebook.GetCurrentPage() {
		case 0:
			principal.savePersonContent(foreignPopUp.PersonContent)
		case 1:
			principal.saveGroupContent(foreignPopUp.GroupContent)
		}
		foreignPopUp.Win.Close()
	})
}

func (principal *Principal) populateFromExistingDB(database *model.Database) {
	rows, err := database.Database.Query("SELECT performers.name, albums.name, rolas.path, rolas.title, rolas.genre, rolas.id_rola FROM rolas INNER JOIN performers ON performers.id_performer = rolas.id_performer INNER JOIN albums ON albums.id_album = rolas.id_album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var performer string
		var album string
		var path string
		var title string
		var genre string
		var id int64
		err = rows.Scan(&performer, &album, &path, &title, &genre, &id)
		if err != nil {
			log.Fatal(err)
		}
		if principal.treeview.Rows[id] == nil {
			glib.IdleAdd(principal.treeview.addRowStruct, &RowInfo{title, performer, album, genre, path, true, id})
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("could not scan row:", err)
	}
}

func (principal *Principal) populateOnTheFly(miner *model.Miner) {
	for rola := range miner.TrackList {
		glib.IdleAdd(principal.treeview.addRowFromRola, rola)
		time.Sleep(100 * time.Nanosecond)
	}
	principal.mainWindow.Buttons["populate"].SetSensitive(true)
	principal.treeSel.SetMode(gtk.SELECTION_SINGLE)
}

func (principal *Principal) searchAction(wildcard string) {
	principal.treeSel.UnselectAll()
	principal.treeview.AllInvisible()
	parser := model.GetParser()
	stmt, queryTerms, ok := parser.Parse(wildcard)
	if ok {
		for _, id := range principal.database.QueryCustom(stmt, queryTerms...) {
			iter := principal.treeview.Rows[id]
			principal.treeview.ListStore.SetValue(iter, 5, true)
		}
	} else {
		for _, id := range principal.database.QuerySimple(stmt) {
			iter := principal.treeview.Rows[id]
			principal.treeview.ListStore.SetValue(iter, 5, true)
		}
	}
}

func (principal *Principal) rowActivated() {
	rowValues := principal.rowTextValues()
	if len(rowValues) == 0 {
		return
	}
	rolaID := principal.rowID()
	rolaPopUp := view.EditRolaWindow()
	rola := principal.database.QueryRola(rolaID)
	principal.showRolaContent(rolaPopUp.RolaContent, rola)

	rolaPopUp.SaveB.Connect("clicked", func() {
		principal.saveRolaContent(rolaPopUp.RolaContent, rolaID, rowValues[4])
		glib.IdleAdd(principal.treeview.updateRow, principal.rolaContentToRow(rolaPopUp.RolaContent, rolaID))
		rolaPopUp.Win.Close()
	})
}

func (principal *Principal) saveGroupContent(groupContent *view.GroupContent) {
	newGroupName := view.GetTextEntry(groupContent.GroupNameE)
	newStart := view.GetTextEntry(groupContent.StartE)
	newEnd := view.GetTextEntry(groupContent.EndE)
	principal.saveGroup(newGroupName, newStart, newEnd)
}

func (principal *Principal) savePersonContent(personContent *view.PersonContent) {
	newStageName := view.GetTextEntry(personContent.StageNameE)
	newRealName := view.GetTextEntry(personContent.RealNameE)
	newBirth := view.GetTextEntry(personContent.BirthE)
	newDeath := view.GetTextEntry(personContent.DeathE)
	principal.savePerson(newStageName, newRealName, newBirth, newDeath)
}

func (principal *Principal) saveRolaContent(rolaContent *view.RolaContent, rolaID int64, path string) {
	rola := model.NewRola()
	rola.SetID(rolaID)
	rola.SetPath(path)
	rola.SetTitle(view.GetTextEntry(rolaContent.TitleE))
	rola.SetArtist(view.GetTextEntry(rolaContent.ArtistE))
	rola.SetAlbum(view.GetTextEntry(rolaContent.AlbumE))
	rola.SetGenre(view.GetTextEntry(rolaContent.GenreE))
	if isInt(view.GetTextEntry(rolaContent.TrackE)) {
		newTrack, _ := strconv.Atoi(view.GetTextEntry(rolaContent.TrackE))
		rola.SetTrack(newTrack)
	} else {
		oldRola := principal.database.QueryRola(rolaID)
		rola.SetTrack(oldRola.Track())
	}
	if isInt(view.GetTextEntry(rolaContent.YearE)) {
		newYear, _ := strconv.Atoi(view.GetTextEntry(rolaContent.YearE))
		rola.SetYear(newYear)
	} else {
		oldRola := principal.database.QueryRola(rolaID)
		rola.SetYear(oldRola.Year())
	}
	principal.database.UpdateRola(rola)
}

func (principal *Principal) saveGroup(groupName, start, end string) {
	groupID := principal.database.ExistsGroup(groupName)
	if groupID > 0 {
		principal.database.UpdateGroup(groupName, start, end, groupID)
	} else {
		principal.database.AddGroup(groupName, start, end)
	}
}

func (principal *Principal) savePerson(stageName, realName, birth, death string) {
	personID := principal.database.ExistsPerson(stageName)
	if personID > 0 {
		principal.database.UpdatePerson(stageName, realName, birth, death, personID)
	} else {
		principal.database.AddPerson(stageName, realName, birth, death)
	}
}

func (principal *Principal) rolaContentToRow(content *view.RolaContent, rolaID int64) *model.Rola {
	rola := model.NewRola()
	rola.SetID(rolaID)
	rola.SetTitle(view.GetTextEntry(content.TitleE))
	rola.SetArtist(view.GetTextEntry(content.ArtistE))
	rola.SetAlbum(view.GetTextEntry(content.AlbumE))
	rola.SetGenre(view.GetTextEntry(content.GenreE))
	track, _ := strconv.Atoi(view.GetTextEntry(content.TrackE))
	rola.SetTrack(track)
	year, _ := strconv.Atoi(view.GetTextEntry(content.YearE))
	rola.SetYear(year)
	return rola
}

func (principal *Principal) showPersonContent(content *view.PersonContent, name string) {
	personID := principal.database.ExistsPerson(name)
	stageName, realName, birth, death := principal.database.QueryPerson(personID)
	glib.IdleAdd(content.StageNameE.SetText, stageName)
	glib.IdleAdd(content.RealNameE.SetText, realName)
	glib.IdleAdd(content.BirthE.SetText, birth)
	glib.IdleAdd(content.DeathE.SetText, death)
}

func (principal *Principal) showGroupContent(content *view.GroupContent, name string) {
	groupID := principal.database.ExistsGroup(name)
	groupName, start, end := principal.database.QueryGroup(groupID)
	glib.IdleAdd(content.GroupNameE.SetText, groupName)
	glib.IdleAdd(content.StartE.SetText, start)
	glib.IdleAdd(content.EndE.SetText, end)
}

func (principal *Principal) showRolaContent(content *view.RolaContent, rola *model.Rola) {
	content.ArtistE.SetText(rola.Artist())
	content.AlbumE.SetText(rola.Album())
	content.TitleE.SetText(rola.Title())
	content.GenreE.SetText(rola.Genre())
	content.TrackE.SetText(strconv.Itoa(rola.Track()))
	content.YearE.SetText(strconv.Itoa(rola.Year()))
}

func (principal *Principal) defaultImage(title, artist, album string) {
	pix, _ := gdk.PixbufNewFromFileAtScale("../data/noimage.png", 250, 250, false)
	image, _ := gtk.ImageNewFromPixbuf(pix)
	glib.IdleAdd(principal.attachInfo, &SongInfo{image, title, artist, album})
	glib.IdleAdd(principal.mainWindow.Win.ShowAll)
}

func (principal *Principal) attachInfo(songInfo *SongInfo) {
	previous, err := principal.mainWindow.Grid.GetChildAt(0, 0)
	if err != nil {
		log.Fatal("unable to get child from grid:", err)
	}
	previous.Destroy()
	principal.mainWindow.Grid.Attach(songInfo.image, 0, 0, 1, 1)
	principal.mainWindow.SongInfo[0].SetText("\n\n\n\t" + songInfo.title + "\n\n\n")
	principal.mainWindow.SongInfo[1].SetText("\t" + songInfo.artist + "\n\n\n")
	principal.mainWindow.SongInfo[2].SetText("\t" + songInfo.album)
	principal.mainWindow.Win.ShowAll()
}

func isInt(text string) bool {
	for _, c := range text {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
