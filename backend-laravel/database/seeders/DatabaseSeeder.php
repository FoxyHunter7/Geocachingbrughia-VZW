<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     */
    public function run(): void
    {
        $this->call([
            LanguageSeeder::class,
            EventSeeder::class,
            EventLanguageSeeder::class,
            GeocacheSeeder::class,
            MessageSeeder::class,
            MessageLanguageSeeder::class,
            StaticSitePropertiesSeeder::class,
            StaticSiteContentSeeder::class,
            SocialsSeeder::class
        ]);
    }
}
