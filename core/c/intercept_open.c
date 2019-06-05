#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/syscalls.h>
#include <linux/unistd.h>
#include <linux/proc_fs.h>
#include <asm/uaccess.h>
#include <linux/namei.h>

int flag=0;

#define MAX_TRY 1024;

MODULE_LICENSE ("GPL");

unsigned long *sys_call_table;

asmlinkage long (*original_sys_open) (const char __user * filename, int
		flags, int mode);

asmlinkage int our_fake_open_function(const char __user *filename, int
		flags, int mode)
{
	int error;
	struct nameidata nd,nd_t;
	struct inode *inode,*inode_t;
	mm_segment_t fs;

	error=user_path_walk(filename,&nd);

	if(!error)
	{

		inode=nd.dentry->d_inode;

		/*Have to do this before calling user_path_walk( )
		 *                 from kernel space:*/
		fs=get_fs( );
		set_fs(get_ds( ));

		/*Protect /tmp/test. Change this to whatever file you
		 *                 want to protect*/
		error=user_path_walk("/tmp/test",&nd_t);

		set_fs(fs);

		if(!error)
		{
			inode_t=nd_t.dentry->d_inode;

			if(inode==inode_t)
				return -EACCES;
		}
	}

	return original_sys_open(filename,flags,mode);
}

static int __init my_init (void)
{
	int i=MAX_TRY;
	unsigned long *sys_table;
	sys_table = (unsigned long *)&system_utsname;

	while(i)
	{
		if(sys_table[__NR_read] == (unsigned long)sys_read)
		{
			sys_call_table=sys_table;
			flag=1;
			break;   
		}
		i--;
		sys_table++;

	}

	if(flag)
	{
		original_sys_open =(void * )xchg(&sys_call_table[__NR_open],
				our_fake_open_function);
	}

	return 0;

}

static void my_exit (void)
{
	xchg(&sys_call_table[__NR_open], original_sys_open);
}

module_init(my_init);
module_exit(my_exit);
