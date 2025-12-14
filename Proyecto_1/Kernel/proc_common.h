#ifndef PROC_COMMON_H
#define PROC_COMMON_H

#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched/signal.h>
#include <linux/mm.h>
#include <linux/sysinfo.h>

#define CARNET "202109303"

/* Obtener memoria total/libre/uso en KB */
static inline void get_sys_mem_kb(unsigned long *total_kb,
                                  unsigned long *free_kb,
                                  unsigned long *used_kb)
{
    struct sysinfo i;
    si_meminfo(&i);
    *total_kb = (unsigned long)i.totalram * i.mem_unit / 1024;
    *free_kb  = (unsigned long)i.freeram  * i.mem_unit / 1024;
    *used_kb  = *total_kb - *free_kb;
}

/* Calcular VSZ y RSS en KB */
static inline void calc_vsz_rss_kb(struct task_struct *task,
                                   unsigned long *vsz_kb,
                                   unsigned long *rss_kb)
{
    if (task->mm) {
        *vsz_kb = task->mm->total_vm * PAGE_SIZE / 1024;
        *rss_kb = get_mm_rss(task->mm) * PAGE_SIZE / 1024;
    } else {
        *vsz_kb = 0;
        *rss_kb = 0;
    }
}

/* %Mem = RSS / TotalRAM */
static inline int calc_mem_percent(unsigned long rss_kb,
                                   unsigned long total_kb)
{
    if (total_kb == 0) return 0;
    return (int)((rss_kb * 100UL) / total_kb);
}

#endif /* PROC_COMMON_H */
