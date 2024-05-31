<?php

namespace Database\Seeders;

use App\Models\EventLanguage;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class EventLanguageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/eventsLanguage.csv'), ';', '/');

        $model = new EventLanguage();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
