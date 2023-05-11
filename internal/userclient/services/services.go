package services

type StorageProvider interface {
	CardStorageProviver
	TextStorageProviver
	BinaryStorageProviver
	LoginStorageProviver
	UsersStorageProviver
}
