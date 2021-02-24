import {MigrationInterface, QueryRunner} from "typeorm";

export class Initial1614018203765 implements MigrationInterface {
    name = 'Initial1614018203765'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "guilds" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "name" character varying NOT NULL, "ownerId" character varying NOT NULL, "icon" character varying, "inviteLinks" text NOT NULL DEFAULT '[]', "membersId" character varying, "membersUserId" character varying, "membersGuildId" character varying, CONSTRAINT "PK_e7e7f2a51bd6d96a9ac2aa560f9" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE INDEX "IDX_6de265be5805c7145f8bf8380d" ON "guilds" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "members" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "userId" character varying NOT NULL, "guildId" character varying NOT NULL, "nickname" character varying, "color" character varying, "lastSeen" TIMESTAMP NOT NULL DEFAULT now(), CONSTRAINT "PK_92a577ea23fc211edab40817a8a" PRIMARY KEY ("id", "userId", "guildId"))`);
        await queryRunner.query(`CREATE INDEX "IDX_12d2c4a4f7ac9472f51610b2f7" ON "members" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "pcmembers" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "userId" character varying NOT NULL, "channelId" character varying NOT NULL, CONSTRAINT "PK_339ffcc6882cccccbb0c55ed6ec" PRIMARY KEY ("id", "userId", "channelId"))`);
        await queryRunner.query(`CREATE INDEX "IDX_7dc754036c1a1b45c1ea4c9121" ON "pcmembers" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "channels" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "name" character varying NOT NULL, "isPublic" boolean NOT NULL DEFAULT true, "dm" boolean NOT NULL DEFAULT false, "lastActivity" TIMESTAMP NOT NULL DEFAULT now(), "guildId" character varying, CONSTRAINT "PK_bc603823f3f741359c2339389f9" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE INDEX "IDX_6480a2e715dda44c2e44826712" ON "channels" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "users" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "username" character varying NOT NULL, "email" character varying NOT NULL, "password" text NOT NULL, "image" text, "isOnline" boolean NOT NULL DEFAULT true, CONSTRAINT "UQ_97672ac88f789774dd47f7c8be3" UNIQUE ("email"), CONSTRAINT "PK_a3ffb1c0c8416b9fc6f907b7433" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE INDEX "IDX_204e9b624861ff4a5b26819210" ON "users" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "bans" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "userId" character varying NOT NULL, "guildId" character varying NOT NULL, CONSTRAINT "PK_8cd40663064e3c087c3e0a6e5a7" PRIMARY KEY ("id", "userId", "guildId"))`);
        await queryRunner.query(`CREATE INDEX "IDX_e37d75f734495b133c132ccb6d" ON "bans" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "dm_members" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "userId" character varying NOT NULL, "channelId" character varying NOT NULL, "isOpen" boolean NOT NULL DEFAULT false, CONSTRAINT "PK_498dbce2de0acbe37c06bd7ff4a" PRIMARY KEY ("id", "userId", "channelId"))`);
        await queryRunner.query(`CREATE INDEX "IDX_ea7bd60fcb572fcf35e9839095" ON "dm_members" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "messages" ("id" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "updatedAt" TIMESTAMP NOT NULL DEFAULT now(), "text" text, "url" text, "filetype" character varying(50), "channelId" character varying, "userId" character varying, CONSTRAINT "PK_18325f38ae6de43878487eff986" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE INDEX "IDX_6ce6acdb0801254590f8a78c08" ON "messages" ("createdAt") `);
        await queryRunner.query(`CREATE TABLE "channel_member" ("channels" character varying NOT NULL, "users" character varying NOT NULL, CONSTRAINT "PK_7274d6926ac49f4d6e3a5626450" PRIMARY KEY ("channels", "users"))`);
        await queryRunner.query(`CREATE INDEX "IDX_efbb9e288fcf35e991ac3696fa" ON "channel_member" ("channels") `);
        await queryRunner.query(`CREATE INDEX "IDX_b80e29716e7e08189d18e723ae" ON "channel_member" ("users") `);
        await queryRunner.query(`CREATE TABLE "friends" ("user" character varying NOT NULL, "friend" character varying NOT NULL, CONSTRAINT "PK_266baf1ac80ab879fc8d4896376" PRIMARY KEY ("user", "friend"))`);
        await queryRunner.query(`CREATE INDEX "IDX_045d2990ee004cfc1401bfb4cf" ON "friends" ("user") `);
        await queryRunner.query(`CREATE INDEX "IDX_10390c9cb4aa6c2364a35eb1ef" ON "friends" ("friend") `);
        await queryRunner.query(`CREATE TABLE "friends_request" ("senderId" character varying NOT NULL, "receiverId" character varying NOT NULL, CONSTRAINT "PK_4e0f09344a7acf3070bcf55d905" PRIMARY KEY ("senderId", "receiverId"))`);
        await queryRunner.query(`CREATE INDEX "IDX_9ddad72959cff6425d0e63a3b4" ON "friends_request" ("senderId") `);
        await queryRunner.query(`CREATE INDEX "IDX_6064cfc70c6df6b01b3469ff67" ON "friends_request" ("receiverId") `);
        await queryRunner.query(`ALTER TABLE "guilds" ADD CONSTRAINT "FK_5097cc50286d7349c8c8c257a3b" FOREIGN KEY ("membersId", "membersUserId", "membersGuildId") REFERENCES "members"("id","userId","guildId") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "FK_16f7ae247a7cf9894db7f23df8e" FOREIGN KEY ("guildId") REFERENCES "guilds"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "messages" ADD CONSTRAINT "FK_fad0fd6def6fa89f66dcf5aaca5" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "messages" ADD CONSTRAINT "FK_4838cd4fc48a6ff2d4aa01aa646" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channel_member" ADD CONSTRAINT "FK_efbb9e288fcf35e991ac3696faa" FOREIGN KEY ("channels") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channel_member" ADD CONSTRAINT "FK_b80e29716e7e08189d18e723ae1" FOREIGN KEY ("users") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "friends" ADD CONSTRAINT "FK_045d2990ee004cfc1401bfb4cf6" FOREIGN KEY ("user") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "friends" ADD CONSTRAINT "FK_10390c9cb4aa6c2364a35eb1ef2" FOREIGN KEY ("friend") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "friends_request" ADD CONSTRAINT "FK_9ddad72959cff6425d0e63a3b41" FOREIGN KEY ("senderId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "friends_request" ADD CONSTRAINT "FK_6064cfc70c6df6b01b3469ff67d" FOREIGN KEY ("receiverId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "friends_request" DROP CONSTRAINT "FK_6064cfc70c6df6b01b3469ff67d"`);
        await queryRunner.query(`ALTER TABLE "friends_request" DROP CONSTRAINT "FK_9ddad72959cff6425d0e63a3b41"`);
        await queryRunner.query(`ALTER TABLE "friends" DROP CONSTRAINT "FK_10390c9cb4aa6c2364a35eb1ef2"`);
        await queryRunner.query(`ALTER TABLE "friends" DROP CONSTRAINT "FK_045d2990ee004cfc1401bfb4cf6"`);
        await queryRunner.query(`ALTER TABLE "channel_member" DROP CONSTRAINT "FK_b80e29716e7e08189d18e723ae1"`);
        await queryRunner.query(`ALTER TABLE "channel_member" DROP CONSTRAINT "FK_efbb9e288fcf35e991ac3696faa"`);
        await queryRunner.query(`ALTER TABLE "messages" DROP CONSTRAINT "FK_4838cd4fc48a6ff2d4aa01aa646"`);
        await queryRunner.query(`ALTER TABLE "messages" DROP CONSTRAINT "FK_fad0fd6def6fa89f66dcf5aaca5"`);
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "FK_16f7ae247a7cf9894db7f23df8e"`);
        await queryRunner.query(`ALTER TABLE "guilds" DROP CONSTRAINT "FK_5097cc50286d7349c8c8c257a3b"`);
        await queryRunner.query(`DROP INDEX "IDX_6064cfc70c6df6b01b3469ff67"`);
        await queryRunner.query(`DROP INDEX "IDX_9ddad72959cff6425d0e63a3b4"`);
        await queryRunner.query(`DROP TABLE "friends_request"`);
        await queryRunner.query(`DROP INDEX "IDX_10390c9cb4aa6c2364a35eb1ef"`);
        await queryRunner.query(`DROP INDEX "IDX_045d2990ee004cfc1401bfb4cf"`);
        await queryRunner.query(`DROP TABLE "friends"`);
        await queryRunner.query(`DROP INDEX "IDX_b80e29716e7e08189d18e723ae"`);
        await queryRunner.query(`DROP INDEX "IDX_efbb9e288fcf35e991ac3696fa"`);
        await queryRunner.query(`DROP TABLE "channel_member"`);
        await queryRunner.query(`DROP INDEX "IDX_6ce6acdb0801254590f8a78c08"`);
        await queryRunner.query(`DROP TABLE "messages"`);
        await queryRunner.query(`DROP INDEX "IDX_ea7bd60fcb572fcf35e9839095"`);
        await queryRunner.query(`DROP TABLE "dm_members"`);
        await queryRunner.query(`DROP INDEX "IDX_e37d75f734495b133c132ccb6d"`);
        await queryRunner.query(`DROP TABLE "bans"`);
        await queryRunner.query(`DROP INDEX "IDX_204e9b624861ff4a5b26819210"`);
        await queryRunner.query(`DROP TABLE "users"`);
        await queryRunner.query(`DROP INDEX "IDX_6480a2e715dda44c2e44826712"`);
        await queryRunner.query(`DROP TABLE "channels"`);
        await queryRunner.query(`DROP INDEX "IDX_7dc754036c1a1b45c1ea4c9121"`);
        await queryRunner.query(`DROP TABLE "pcmembers"`);
        await queryRunner.query(`DROP INDEX "IDX_12d2c4a4f7ac9472f51610b2f7"`);
        await queryRunner.query(`DROP TABLE "members"`);
        await queryRunner.query(`DROP INDEX "IDX_6de265be5805c7145f8bf8380d"`);
        await queryRunner.query(`DROP TABLE "guilds"`);
    }

}
