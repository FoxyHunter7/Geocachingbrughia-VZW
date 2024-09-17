<?php

namespace Database\Seeders;

use App\Models\Language;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class LanguageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('languages')->truncate();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/languages.csv'), ';');

        $model = new Language();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
