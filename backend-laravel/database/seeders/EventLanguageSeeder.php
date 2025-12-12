<?php

namespace Database\Seeders;

use App\Models\EventLanguage;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class EventLanguageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('event_languages')->delete();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/eventsLanguage.csv'), ';', '/');

        $model = new EventLanguage();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
