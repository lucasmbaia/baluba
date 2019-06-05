#include <linux/init.h>	  //usado para marcar funcoes de inicializao e fim __init __exit
#include <linux/module.h> //pacote principal de headers para carregar LKMS no kernel
#include <linux/device.h>
#include <linux/kernel.h>
#include <linux/fs.h>
#include <linux/uaccess.h>
#define	DEVICE_NAME "balubafs"
#define	CLASS_NAME  "baluba"

#define AUFS_MAGIC_NUMBER 0x10032013

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Lucas Martins Baia");
MODULE_DESCRIPTION("Simple File System to baluba");
MODULE_VERSION("0.1");

static int baluba_open(struct inode *, struct file *f);
static ssize_t baluba_read(struct file *f, char __user *buf, size_t len, loff_t *ppos);
static ssize_t baluba_write(struct file *f, const char *buf, size_t len, loff_t *ppos);

static void balubafs_put_super(struct super_block *sb) {
	pr_debug("balubafs super block destroyed\n");
}

static struct super_operations const balubafs_super_ops =  {
	.put_super = balubafs_put_super,
};

static int balubafs_fill_sb(struct super_block *sb, void *data, int silent) {
	struct inode *root = NULL;

	sb->s_magic = AUFS_MAGIC_NUMBER;
	sb->s_op = &balubafs_super_ops;

	root = new_inode(sb);
	if (!root) {
		pr_err("inode allocation failed\n");
		return -ENOMEM;
	}

	root->i_ino = 0;
	root->i_sb = sb;
	root->i_atime = root->i_mtime = root->i_ctime = CURRENT_TIME;
	inode_init_owner(root, NULL, S_IFDIR);

	sb->s_root = d_make_root(root);
	if (!sb->s_root) {
		pr_err("root creation failed\n");
		return -ENOMEM;
	}

	return 0;
}

static struct dentry *balubafs_mount(struct file_system_type *type, int flags, char const *dev, void *data) {
	struct dentry *const entry = mount_bdev(type, flags, dev, data, balubafs_fill_sb);

	if (IS_ERR(entry)) {
		pr_err("auts mounting failed\n");
	} else {
		pr_debug("aufs mounted\n");
	}

	return entry;
}

static struct file_operations baluba_file_operations = {
	.owner = THIS_MODULE,
	.open = baluba_open,
	.read = baluba_read,
	.write = baluba_write,
};

static struct file_system_type balubafs_fs_type = {
	.owner = THIS_MODULE,
	.name = DEVICE_NAME,
	.mount = balubafs_mount,
	.kill_sb = kill_block_super,
	.fs_flags = FS_REQUIRES_DEV,
};

static int balubafs_init(void) {
	int ret;

	ret = register_filesystem(&balubafs_fs_type);
	if (likely(ret == 0)) {
		printk(KERN_INFO "Sucessfully registered balubafs\n");
	} else {
		printk(KERN_ERR "Failed to register balubafs. Error:[%d]", ret);
	}

	return ret;
}

static void __exit balubafs_exit(void) {
	int ret;

	ret = unregister_filesystem(&balubafs_fs_type);

	if (likely(ret == 0)) {
		printk(KERN_INFO "Sucessfully unregistered balubafs\n");
	} else {
		printk(KERN_ERR "Failed to unregister balubafs. Error:[%d]", ret);
	}
}

static int baluba_open(struct inode *inodep, struct file *f) {
	return 0;
}

static ssize_t baluba_read(struct file *f, char __user *buf, size_t len, loff_t *ppos) {
	return len;
}

static ssize_t baluba_write(struct file *f, const char *buf, size_t len, loff_t *ppos) {
	return len;
}

module_init(balubafs_init);
module_exit(balubafs_exit);
