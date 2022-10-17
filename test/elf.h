#ifndef _ELF_H_
#define _ELF_H_

#include "util/types.h"
#include "process.h"

#define MAX_CMDLINE_ARGS 64
#define MAX_SYMOBL_NUMS 128
#define MAX_STRING_NUMS 4096
// elf symbol structure
typedef struct
{
  uint32    st_name;   /* Symbol name (string tbl index) */
  unsigned char st_info;   /* Symbol type and binding */
  unsigned char st_other;  /* Symbol visibility */
  uint16 st_shndx;  /* Section index */
  uint64    st_value;  /* Symbol value */
  uint64   st_size;   /* Symbol size */
}elf_sym;

/* This info is needed when parsing the symbol table */
/* 这些是st_info 的数值类型 */
#define STB_LOCAL  0
#define STB_GLOBAL 1
#define STB_WEAK   2

#define STT_NOTYPE  0
#define STT_OBJECT  1
#define STT_FUNC    18
#define STT_SECTION 3
#define STT_FILE    4
#define STT_COMMON  5
#define STT_TLS     6

// elf header structure
typedef struct elf_header_t {
  uint32 magic;
  uint8 elf[12];
  uint16 type;      /* Object file type */
  uint16 machine;   /* Architecture */
  uint32 version;   /* Object file version */
  uint64 entry;     /* Entry point virtual address */
  uint64 phoff;     /* Program header table file offset */
  uint64 shoff;     /* Section header table file offset */
  uint32 flags;     /* Processor-specific flags */
  uint16 ehsize;    /* ELF header size in bytes */
  uint16 phentsize; /* Program header table entry size */
  uint16 phnum;     /* Program header table entry count */
  uint16 shentsize; /* Section header table entry size */
  uint16 shnum;     /* Section header table entry count */
  uint16 shstrndx;  /* Section header string table index */
} elf_header;

// elf section header structure 
// 参考自
// https://mudongliang.github.io/2015/05/19/elf-section-header-data-structure.html
typedef struct
{
  uint32 sh_name;      /* Section name (string tbl index) */
  uint32 sh_type;      /* Section type */
  uint64 sh_flags;     /* Section flags */
  uint64 sh_addr;      /* Section virtual addr at execution */
  uint64 sh_offset;    /* Section file offset */
  uint64 sh_size;      /* Section size in bytes */
  uint32 sh_link;      /* Link to another section */
  uint32 sh_info;      /* Additional section information */
  uint64 sh_addralign; /* Section alignment */
  uint64 sh_entsize;   /* Entry size if section holds table */
} elf_section_header;
// Program segment header.
typedef struct elf_prog_header_t {
  uint32 type;   /* Segment type */
  uint32 flags;  /* Segment flags */
  uint64 off;    /* Segment file offset */
  uint64 vaddr;  /* Segment virtual address */
  uint64 paddr;  /* Segment physical address */
  uint64 filesz; /* Segment size in file */
  uint64 memsz;  /* Segment size in memory */
  uint64 align;  /* Segment alignment */
} elf_prog_header;

//  The Segment type
// 参考自
// /usr/include/elf.h
#define SHT_NULL	0
#define SHT_PROGBITS	1
#define SHT_SYMTAB	2
#define SHT_STRTAB	3



#define ELF_MAGIC 0x464C457FU  // "\x7FELF" in little endian
#define ELF_PROG_LOAD 1

typedef enum elf_status_t {
  EL_OK = 0,

  EL_EIO,
  EL_ENOMEM,
  EL_NOTELF,
  EL_ERR,

} elf_status;

typedef struct elf_ctx_t {
  void *info;
  elf_header ehdr;
  elf_sym sym_table[MAX_SYMOBL_NUMS];
  char str_table[MAX_STRING_NUMS];
  uint64 sys_nums;
} elf_ctx;
elf_status elf_init(elf_ctx *ctx, void *info);
elf_status elf_load(elf_ctx *ctx);
elf_status elf_load_sym(elf_ctx *ctx);

void load_bincode_from_host_elf(process *p);

#endif
