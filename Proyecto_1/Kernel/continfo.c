#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched/signal.h>
#include <linux/mm.h>
#include <linux/uaccess.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Marvin Paz");
MODULE_DESCRIPTION("Contenedores info module");

static bool is_container_process(const char *comm) {
    // Ajuste: detecta procesos típicos de tus contenedores
    if (strcmp(comm, "sleep") == 0) return true;
    if (strcmp(comm, "python3") == 0) return true;
    return false;
}

static int continfo_show(struct seq_file *m, void *v) {
    struct task_struct *task;
    unsigned long totalram = totalram_pages() << (PAGE_SHIFT - 10);
    unsigned long freeram = global_zone_page_state(NR_FREE_PAGES) << (PAGE_SHIFT - 10);
    unsigned long usedram = totalram - freeram;

    seq_printf(m,
        "{\n  \"memory\": {\"total_kb\": %lu, \"free_kb\": %lu, \"used_kb\": %lu},\n",
        totalram, freeram, usedram);

    seq_puts(m, "  \"processes\": [\n");

    bool first = true;
    for_each_process(task) {
        if (task->flags & PF_KTHREAD) continue; // ignora hilos del kernel
        if (!is_container_process(task->comm)) continue;

        unsigned long rss = task->mm ? get_mm_rss(task->mm) << (PAGE_SHIFT - 10) : 0;
        unsigned long vsz = task->mm ? task->mm->total_vm << (PAGE_SHIFT - 10) : 0;
        int mem_percent = totalram ? (rss * 100) / totalram : 0;

        if (!first) seq_puts(m, ",\n");
        first = false;

        seq_printf(m,
            "    {\"pid\": %d, \"name\": \"%s\", \"vsz_kb\": %lu, \"rss_kb\": %lu, \"mem_percent\": %d}",
            task->pid, task->comm, vsz, rss, mem_percent);
    }

    seq_puts(m, "\n  ]\n}\n");
    return 0;
}

static int continfo_open(struct inode *inode, struct file *file) {
    return single_open(file, continfo_show, NULL);
}

static const struct proc_ops continfo_fops = {
    .proc_open    = continfo_open,
    .proc_read    = seq_read,
    .proc_lseek   = seq_lseek,
    .proc_release = single_release,
};

static int __init continfo_init(void) {
    proc_create("continfo_so1_202109303", 0, NULL, &continfo_fops);
    pr_info("continfo module loaded\n");
    return 0;
}

static void __exit continfo_exit(void) {
    remove_proc_entry("continfo_so1_202109303", NULL);
    pr_info("continfo module unloaded\n");
}

module_init(continfo_init);
module_exit(continfo_exit);
