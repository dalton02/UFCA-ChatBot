-- CreateTable
CREATE TABLE "Usuario" (
    "Id" SERIAL NOT NULL,
    "Tipo" TEXT NOT NULL,
    "Login" TEXT NOT NULL,
    "Senha" TEXT,

    CONSTRAINT "Usuario_pkey" PRIMARY KEY ("Id")
);
