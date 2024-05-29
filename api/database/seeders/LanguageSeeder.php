<?php

namespace Database\Seeders;

use App\Models\Language;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class LanguageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/languages.csv'), ';');

        $model = new Language();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
