#include "proc_common.h"

#define PROC_NAME "sysinfo_so1_" CARNET

static int proc_show(struct seq_file *m, void *v)
{
    unsigned long total_kb, free_kb, used_kb;
    get_sys_mem_kb(&total_kb, &free_kb, &used_kb);

    seq_printf(m, "{\n  \"memory\": {\"total_kb\": %lu, \"free_kb\": %lu, \"used_kb\": %lu},\n",
               total_kb, free_kb, used_kb);
    seq_puts(m, "  \"processes\": [\n");

    bool first = true;
    struct task_struct *task;

    rcu_read_lock();
    for_each_process(task) {
        if (task->flags & PF_KTHREAD) continue;

        unsigned long vsz, rss;
        calc_vsz_rss_kb(task, &vsz, &rss);
        int mem_pct = calc_mem_percent(rss, total_kb);

        char comm[TASK_COMM_LEN];
        get_task_comm(comm, task);

        if (!first) seq_puts(m, ",\n");
        first = false;

        seq_printf(m,
            "    {\"pid\": %d, \"name\": \"%s\", \"vsz_kb\": %lu, \"rss_kb\": %lu, \"mem_percent\": %d}",
            task->pid, comm, vsz, rss, mem_pct);
    }
    rcu_read_unlock();

    seq_puts(m, "\n  ]\n}\n");
    return 0;
}

static int proc_open(struct inode *inode, struct file *file)
{
    return single_open(file, proc_show, NULL);
}

static const struct proc_ops pops = {
    .proc_open    = proc_open,
    .proc_read    = seq_read,
    .proc_lseek   = seq_lseek,
    .proc_release = single_release,
};

static int __init sysinfo_init(void)
{
    if (!proc_create(PROC_NAME, 0444, NULL, &pops)) {
        pr_err("Failed to create /proc/%s\n", PROC_NAME);
        return -ENOMEM;
    }
    pr_info("sysinfo module loaded\n");
    return 0;
}

static void __exit sysinfo_exit(void)
{
    remove_proc_entry(PROC_NAME, NULL);
    pr_info("sysinfo module unloaded\n");
}

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Marvin");
MODULE_DESCRIPTION("System processes info in /proc (JSON)");

module_init(sysinfo_init);
module_exit(sysinfo_exit);
